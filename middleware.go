package main

import (
	"net/http"
	"time"

	"github.com/rs/xlog"
)

func accessLoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		logger := xlog.FromRequest(r)
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		logger.Infof(
			"access-log",
			xlog.F{
				"req_time":      t2.Sub(t1),
				"req_time_nsec": t2.Sub(t1).Nanoseconds(),
			},
		)
	}
	return http.HandlerFunc(fn)
}
