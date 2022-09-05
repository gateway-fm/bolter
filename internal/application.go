package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gateway-fm/scriptorium/logger"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/hcl/v2/hclsimple"
	vegeta "github.com/tsenart/vegeta/v12/lib"

	"bolter/config"
	"bolter/internal/models"
)

type Bolter struct {
	Request   models.JsonBody
	Vegeta    models.Vegeta
	Logger    models.Logger
	BolterCfg config.BolterCfg
}

func LoadBolter() error {

	bolt := &Bolter{}
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

	targeter, err := bolt.PrepareTargeter()
	if err != nil {
		return fmt.Errorf("PrepareRate failed %w", err)
	}

	attacker := bolt.PrepareAttacker()

	var metrics vegeta.Metrics
	var result *models.ResultBody
	var wg sync.WaitGroup

	for res := range attacker.Attack(*targeter, rate, *duration, "") {
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

func (b *Bolter) NewBody() ([]byte, error) {
	cfg, err := b.DecodeConfig()
	if err != nil {
		return nil, fmt.Errorf("decoding failed: %w", err)
	}
	b.Request.Method = cfg.Request.Method
	b.Request.Id = cfg.Request.Id
	b.Request.Jsonrpc = cfg.Request.Jsonrpc
	b.Request.Params = cfg.Request.Parameters
	body, err := json.Marshal(b.Request)
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

// PrepareTargeter loading new vegeta.Targeter
func (b *Bolter) PrepareTargeter() (*vegeta.Targeter, error) {
	cfg, err := b.DecodeConfig()
	if err != nil {
		return nil, fmt.Errorf("decoding failed: %w", err)
	}
	body, err := b.NewBody()
	if err != nil {
		return nil, fmt.Errorf("failed getting new body: %w", err)
	}
	b.Vegeta.Target.URL, b.Vegeta.Target.Method = cfg.Vegeta.Url, cfg.Vegeta.Method
	b.Vegeta.Target.Body = body
	if !cfg.Vegeta.IsPublic {
		auth := cfg.Vegeta.Header.Auth + " " + cfg.Vegeta.Header.Bearer
		header := b.Vegeta.Target.Header
		header.Add("Authorization", auth)
	}
	b.Vegeta.Targeter = vegeta.NewStaticTargeter(b.Vegeta.Target)
	return &b.Vegeta.Targeter, nil
}

// PrepareRate loading new vegeta.Rate
func (b *Bolter) PrepareRate() (*vegeta.Rate, error) {
	cfg, err := b.DecodeConfig()
	if err != nil {
		return nil, fmt.Errorf("decoding failed: %w", err)
	}
	b.Vegeta.Rate.Freq = cfg.Vegeta.Rate
	b.Vegeta.Rate.Per = time.Second
	return &b.Vegeta.Rate, nil
}

// PrepareDuration setting Duration for our attacker
func (b *Bolter) PrepareDuration() (*time.Duration, error) {
	cfg, err := b.DecodeConfig()
	if err != nil {
		return nil, fmt.Errorf("decoding failed: %w", err)
	}
	duration := time.Duration(cfg.Vegeta.Duration) * time.Second
	b.Vegeta.Duration = duration
	return &b.Vegeta.Duration, nil
}

func (b *Bolter) CreateLoggerTextFile() (*os.File, error) {
	cfg, err := b.DecodeConfig()
	if err != nil {
		return nil, fmt.Errorf("decoding failed: %w", err)
	}
	b.Logger.LoggerType = cfg.Logger.LoggerType

	if b.Logger.LoggerType == 0 {
		file, err := b.Logger.TextFile.CreateNewTextFile(cfg.Logger.FileName)
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
	cfg, err := b.DecodeConfig()
	if err != nil {
		return fmt.Errorf("decoding failed: %w", err)
	}
	b.Logger.LoggerType = cfg.Logger.LoggerType
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
