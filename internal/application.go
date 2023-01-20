package internal

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/goccy/go-json"

	"github.com/hashicorp/hcl/v2/hclsimple"

	vegeta "github.com/tsenart/vegeta/v12/lib"

	"github.com/gateway-fm/scriptorium/logger"

	"bolter/config"
	"bolter/internal/models"
)

type Bolter struct {
	Request   []models.JsonBody
	Vegeta    models.Vegeta
	Logger    models.Logger
	BolterCfg config.BolterCfg
}

func LoadBolter() error {

	bolt := &Bolter{}
	err := bolt.initCfg()
	if err != nil {
		return fmt.Errorf("unable to init bolter's cfg: %w", err)
	}
	file, err := bolt.CreateLoggerTextFile()
	if err != nil {
		return fmt.Errorf("CreateLoggerTextFile failed %w", err)
	}
	rate, err := bolt.PrepareRate()
	if err != nil {
		return fmt.Errorf("PrepareRate failed %w", err)
	}
	duration, err := bolt.PrepareDuration()
	if err != nil {
		return fmt.Errorf("PrepareDuration failed %w", err)
	}

	attacker := bolt.PrepareAttacker()
	var metrics vegeta.Metrics
	var wg sync.WaitGroup
	trgts, err := bolt.PrepareTargeter()
	if err != nil {
		return fmt.Errorf("PrepareTargeter failed %w", err)
	}
	//result := make([]*models.ResultBody, len(trgts)

	var result *models.ResultBody
	//for i := range trgts {
	//	for res := range attacker.Attack(trgts[i], rate, *duration, "") {
	for res := range attacker.Attack(trgts[0], rate, *duration, "") {
		wg.Add(1)
		go func() {
			defer wg.Done()
			metrics.Add(res)
			err = json.Unmarshal(res.Body, &result)
			if err != nil {
				logger.Log().Error("An error occurred", zap.Error(err))
			}
			err = bolt.PrepareLogger(result.Result, file)
			if err != nil {
				logger.Log().Error("An error occurred", zap.Error(err))
			}
		}()
		wg.Wait()
	}
	//}
	metrics.Close()

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	return nil
}

func (b *Bolter) DecodeConfig() (*config.BolterCfg, error) {
	err := hclsimple.DecodeFile("config/config.hcl", nil, &b.BolterCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	return &b.BolterCfg, nil
}
func (b *Bolter) initCfg() error {
	_, err := b.DecodeConfig()
	if err != nil {
		return fmt.Errorf("decoding failed: %w", err)
	}
	return nil
}

func (b *Bolter) BuildRequests() []models.JsonBody {
	confReqs := b.BolterCfg.Requests
	b.Request = make([]models.JsonBody, len(confReqs))
	for i := range confReqs {
		b.Request[i].Method = confReqs[i].Request.Method
		b.Request[i].Id = confReqs[i].Request.Id
		b.Request[i].Jsonrpc = confReqs[i].Request.Jsonrpc
		b.Request[i].Params = confReqs[i].Request.Parameters
	}
	return b.Request
}
func (b *Bolter) NewBodies(i int) ([]byte, error) {
	reqs := b.BuildRequests()
	body, err := json.Marshal(reqs[i])
	if err != nil {
		return nil, fmt.Errorf("marshalling failed: %w", err)
	}
	return body, nil
}

// PrepareAttacker loading new vegeta.Attacker
func (b *Bolter) PrepareAttacker() *vegeta.Attacker {
	b.Vegeta.Attacker = vegeta.NewAttacker()
	return b.Vegeta.Attacker
}

func (b *Bolter) PrepareTarget() ([]vegeta.Target, error) {
	le := len(b.BolterCfg.Requests)
	var err error
	trg := make([]vegeta.Target, le)
	for i := range trg {
		b.Vegeta.Target.Body, err = b.NewBodies(i)
		if err != nil {
			return nil, fmt.Errorf("failed while NewBodies: %w", err)
		}
		b.Vegeta.Target.URL, b.Vegeta.Target.Method = b.BolterCfg.Vegeta.Url, b.BolterCfg.Vegeta.Method
		if !b.BolterCfg.Vegeta.IsPublic {
			auth := b.BolterCfg.Vegeta.Header.Auth + " " + b.BolterCfg.Vegeta.Header.Bearer
			b.Vegeta.Target.Header = make(http.Header)
			b.Vegeta.Target.Header.Add("Authorization", auth)
		}
		trg = append(trg, b.Vegeta.Target)

	}
	trg = append(trg[le:le], trg[le:]...)

	return trg, err
}

// PrepareTargeter loading new vegeta.Targeter
func (b *Bolter) PrepareTargeter() ([]vegeta.Targeter, error) {
	le := len(b.BolterCfg.Requests)
	trgr := make([]vegeta.Targeter, le)
	prep, err := b.PrepareTarget()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare Target:%w", err)
	}
	for i := range trgr {
		b.Vegeta.Targeter = vegeta.NewStaticTargeter(prep[i])
		trgr = append(trgr, b.Vegeta.Targeter)
	}
	trgr = append(trgr[le:le], trgr[le:]...)

	return trgr, nil
}

// PrepareRate loading new vegeta.Rate
func (b *Bolter) PrepareRate() (*vegeta.Rate, error) {
	b.Vegeta.Rate.Freq = b.BolterCfg.Vegeta.Rate
	b.Vegeta.Rate.Per = time.Second
	return &b.Vegeta.Rate, nil
}

// PrepareDuration setting Duration for our attacker
func (b *Bolter) PrepareDuration() (*time.Duration, error) {
	duration := time.Duration(b.BolterCfg.Vegeta.Duration) * time.Second
	b.Vegeta.Duration = duration
	return &b.Vegeta.Duration, nil
}

func (b *Bolter) CreateLoggerTextFile() (*os.File, error) {
	b.Logger.LoggerType = b.BolterCfg.Logger.LoggerType

	if b.Logger.LoggerType == 0 {
		file, err := b.Logger.TextFile.CreateNewTextFile(b.BolterCfg.Logger.FileName)
		if err != nil {
			return nil, fmt.Errorf("creating new file failed: %w", err)
		}
		return file, nil
	} else {
		return nil, nil
	}
}

// PrepareLogger setting Duration for our attacker
func (b *Bolter) PrepareLogger(data interface{}, file ...*os.File) error {
	var err error
	b.Logger.LoggerType = b.BolterCfg.Logger.LoggerType
	switch b.Logger.LoggerType {
	case 0:
		b.Logger.TextFile.LogWithTextFile(data, file[0])
		if err != nil {
			return fmt.Errorf("couldn't log with text file: %w", err)
		}
	case 1:
		b.Logger.Logrus.LogWithLogrus(data)
	case 2:
		b.Logger.Zap.LogWithZap(data)
	default:
		err = errors.New("logger hasn't been set")
		return err
	}
	return nil
}
