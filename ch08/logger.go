package main

import (
	"os"

	"github.com/rs/xlog"
)

// NewLogConfig creates xlog config
func NewLogConfig(cfg *AppConfig) xlog.Config {
	conf := xlog.Config{
		Level: xlog.LevelDebug,
		Fields: xlog.F{
			"role": "my-service",
		},
		Output: xlog.NewOutputChannel(xlog.NewJSONOutput(os.Stdout)),
	}
	return conf
}
