package main

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/rs/xlog"
)

type ctxKeyType int

const (
	ctxKeyUser ctxKeyType = iota
)

type authUser struct {
	ID   int64
	Name string
}

func faceAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := &authUser{
			ID:   10,
			Name: "achiku",
		}
		ctx := context.WithValue(r.Context(), ctxKeyUser, u)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func getUserFromContext(ctx context.Context) *authUser {
	u, ok := ctx.Value(ctxKeyUser).(*authUser)
	if !ok {
		return &authUser{
			ID:   0,
			Name: "anonymous",
		}
	}
	return u
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := xlog.FromRequest(r)
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		t := t2.Sub(t1)
		logger.Info(
			"accesslog",
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
