package main

import (
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

// AppHandler application handler adaptor
type AppHandler struct {
	h func(http.ResponseWriter, *http.Request) (int, interface{}, error)
}

func (a AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	status, res, err := a.h(w, r)
	if err != nil {
		log.Printf("error: %s", err)
		w.WriteHeader(status)
		encoder.Encode(res)
		return
	}
	w.WriteHeader(status)
	encoder.Encode(res)
	return
}

// Greeting greeting
func (app *App) Greeting(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	res, err := HelloService(r.Context(), "", time.Now())
	if err != nil {
		app.Logger.Printf("error: %s", err)
		e := ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "something went wrong",
		}
		return http.StatusInternalServerError, e, err
	}
	app.Logger.Printf("ok: %v", res)
	return http.StatusOK, res, nil
}

// GreetingWithName greeting with name
func (app *App) GreetingWithName(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	val := mux.Vars(r)
	res, err := HelloService(r.Context(), val["name"], time.Now())
	if err != nil {
		app.Logger.Printf("error: %s", err)
		e := ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "something went wrong",
		}
		return http.StatusInternalServerError, e, err
	}
	app.Logger.Printf("ok: %v", res)
	return http.StatusOK, res, nil
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
	// middleware chain
	// chain := alice.New(
	// 	loggingMiddleware,
	// 	appLoggingMiddleware(app.Logger),
	// )
	// for gorilla/mux
	router := mux.NewRouter()
	r := router.PathPrefix("/api").Subrouter()
	r.Methods("GET").Path("/hello").Handler(
		loggingMiddleware(AppHandler{h: app.Greeting}))
	r.Methods("GET").Path("/hello/staticName").Handler(
		loggingMiddleware(AppHandler{h: app.Greeting}))
	r.Methods("GET").Path("/hello/{name}").Handler(
		loggingMiddleware(AppHandler{h: app.GreetingWithName}))

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
