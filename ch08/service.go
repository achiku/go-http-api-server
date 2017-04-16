package main

import (
	"context"
	"time"
)

// HelloService hello service
func HelloService(ctx context.Context, name string, tm time.Time) (*Greeting, error) {
	if name == "" {
		return &Greeting{
			Message: "hello!",
			Name:    "anonymous",
		}, nil
	} else if name == "8maki" {
		var msg string
		noon := time.Date(tm.Year(), tm.Month(), tm.Day(), 12, 0, 0, 0, time.Local)
		evening := time.Date(tm.Year(), tm.Month(), tm.Day(), 18, 0, 0, 0, time.Local)
		switch {
		case tm.Before(noon):
			msg = "good morning"
		case tm.After(noon) && tm.Before(evening):
			msg = "good afternoon"
		case tm.After(evening):
			msg = "good evening"
		}
		return &Greeting{
			Message: msg,
			Name:    "my boss",
		}, nil
	} else if name == "moqada" {
		return &Greeting{
			Message: "sup",
			Name:    "my man",
		}, nil
	}
	return &Greeting{
		Message: "hello",
		Name:    name,
	}, nil
}
