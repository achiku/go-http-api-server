package main

import (
	"time"

	"github.com/rs/xlog"
)

// Logger logger
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

// App some kind of app
type App struct {
	Logger Logger
}

func main() {
	l := xlog.New(xlog.Config{
		Output: xlog.NewConsoleOutput(),
	})
	ap := App{
		Logger: l,
	}
	ap.Logger.Debug("test log")
	ap.Logger.Debugf("test log: %s", time.Now())
}
