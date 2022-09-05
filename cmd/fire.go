package main

import (
	"fmt"
	"os"

	"github.com/gateway-fm/scriptorium/logger"

	"bolter/cmd/serve"
)

func main() {
	logger.SetLoggerMode("local")
	bolter := serve.Cmd()
	bolter.AddCommand(serve.Cmd())
	if err := bolter.Execute(); err != nil {
		logger.Log().Error(fmt.Errorf("fire has been failed: %w", err).Error())
		os.Exit(1)
	}
}
