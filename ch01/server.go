package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Greeting greeting
type Greeting struct {
	Message string `json:"message"`
	Name    string `json:"name"`
}

// ErrorResponse error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

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

// Greeting greeting
func (app *App) Greeting(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	res, err := HelloService(r.Context(), "", time.Now())
	if err != nil {
		app.Logger.Printf("error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "something went wrong",
		})
		return
	}
	app.Logger.Printf("ok: %v", res)
	w.WriteHeader(http.StatusOK)
	encoder.Encode(res)
	return
}

// GreetingWithName greeting with name
func (app *App) GreetingWithName(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	val := mux.Vars(r)
	res, err := HelloService(r.Context(), val["name"], time.Now())
	if err != nil {
		app.Logger.Printf("error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "something went wrong",
		})
		return
	}
	app.Logger.Printf("ok: %v", res)
	w.WriteHeader(http.StatusOK)
	encoder.Encode(res)
	return
}

// App application
type App struct {
	Host   string
	Name   string
	Logger *log.Logger
}

func main() {
	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	app := App{
		Name:   "my-service",
		Host:   host,
		Logger: log.New(os.Stdout, fmt.Sprintf("[host=%s] ", host), log.LstdFlags),
	}
	// for gorilla/mux
	router := mux.NewRouter()
	r := router.PathPrefix("/api").Subrouter()
	r.Methods("GET").Path("/hello").HandlerFunc(app.Greeting)
	r.Methods("GET").Path("/hello/staticName").HandlerFunc(app.Greeting)
	r.Methods("GET").Path("/hello/{name}").HandlerFunc(app.GreetingWithName)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
