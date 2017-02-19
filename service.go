package main

import (
	"context"
	"fmt"

	"github.com/rs/xlog"
)

// HelloService hello service
func HelloService(ctx context.Context, name string) (string, error) {
	logger := xlog.FromContext(ctx)
	if name == "" {
		logger.Info("no name")
		return "hello!", nil
	} else if name == "8maki" {
		logger.Info("CEO")
		return "hello, Sir!", nil
	} else if name == "moqada" {
		logger.Info("my man")
		return "sup man", nil
	}
	return fmt.Sprintf("hello, %s!", name), nil
}
