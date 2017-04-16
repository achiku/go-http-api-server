package main

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/rs/xlog"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := xlog.FromRequest(r)
		logger.Info("access-log-pre")
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		t := t2.Sub(t1)
		logger.Info(
			"access-log-post",
			xlog.F{
				"req_time":      t.String(),
				"req_time_nsec": t.Nanoseconds(),
			},
		)
	})
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
				log.Printf("panic: %s", err)
				http.Error(w, http.StatusText(
					http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
