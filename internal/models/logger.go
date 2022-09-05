package models

import (
	"fmt"
	lr "github.com/Misnaged/annales/logger"
	zaplog "github.com/gateway-fm/scriptorium/logger"
	"log"
	"os"
)

type Logger struct {
	LoggerType int
	/*
		0: write to text file
		1: write using Logrus
		2: write using ZapLogger
	*/
	TextFile *TextFile
	Logrus   *Logrus
	Zap      *Zap
}
type TextFile struct{}
type Logrus struct{}
type Zap struct{}

func (t *TextFile) CreateNewTextFile(name string) (*os.File, error) {
	file, err := os.Create(name)
	if err != nil {
		return nil, fmt.Errorf("failed to create new file: %w", err)
	}
	return file, nil
}
func (t *TextFile) LogWithTextFile(data interface{}, file *os.File) {
	logger := log.New(file, "", 0)
	logger.Println(data)
}
func (l *Logrus) LogWithLogrus(data interface{}) {
	lr.Log().Infoln(data)
}
func (z *Zap) LogWithZap(data interface{}) {
	zaplog.Log().Info(data.(string))
}
