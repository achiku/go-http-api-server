package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/pkg/errors"
	"github.com/rs/xlog"
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
	logger := xlog.FromRequest(r)
	res, err := HelloService(r.Context(), "", time.Now())
	if err != nil {
		e := ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "something went wrong",
		}
		return http.StatusInternalServerError, e, err
	}
	logger.Debugf("%s %s", res.Name, res.Message)
	return http.StatusOK, res, nil
}

// GreetingWithName greeting with name
func (app *App) GreetingWithName(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	logger := xlog.FromRequest(r)
	name := mux.Vars(r)["name"]
	logger.Debugf("param: %s", name)
	res, err := HelloService(r.Context(), name, time.Now())
	if err != nil {
		e := ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "something went wrong",
		}
		return http.StatusInternalServerError, e, err
	}
	logger.Debugf("%s %s", res.Name, res.Message)
	return http.StatusOK, res, nil
}

// App application
type App struct {
	Host   string
	Name   string
	Config *AppConfig
}

// NewApp creates app
func NewApp(path string) (*App, error) {
	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	cfg, err := NewAppConfig(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load config: %s", path)
	}
	app := &App{
		Name:   "my-service",
		Host:   host,
		Config: cfg,
	}
	return app, nil
}

func main() {
	app, err := NewApp("./devel.toml")
	if err != nil {
		log.Fatal(err)
	}

	// middleware chain
	chain := alice.New(
		recoverMiddleware,
	)
	apiChain := chain.Append(
		xlog.NewHandler(NewLogConfig(app.Config)),
		xlog.MethodHandler("method"),
		xlog.URLHandler("url"),
		xlog.RemoteAddrHandler("ip"),
		xlog.UserAgentHandler("user_agent"),
		xlog.RefererHandler("referer"),
		xlog.RequestIDHandler("req_id", "Request-Id"),
		loggingMiddleware,
	)
	halfLogChain := chain.Append(
		xlog.NewHandler(NewLogConfig(app.Config)),
		loggingMiddleware,
	)
	noLogChain := chain.Append(
		loggingMiddleware,
	)
	// for gorilla/mux
	router := mux.NewRouter()
	r := router.PathPrefix("/api").Subrouter()
	r.Methods("GET").Path("/hello").Handler(apiChain.Then(AppHandler{h: app.Greeting}))
	r.Methods("GET").Path("/hello/nolog").Handler(noLogChain.Then(AppHandler{h: app.Greeting}))
	r.Methods("GET").Path("/hello/halflog").Handler(halfLogChain.Then(AppHandler{h: app.Greeting}))
	r.Methods("GET").Path("/hello/staticName").Handler(apiChain.Then(AppHandler{h: app.Greeting}))
	r.Methods("GET").Path("/hello/{name}").Handler(apiChain.Then(AppHandler{h: app.GreetingWithName}))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", app.Config.ServerPort), router); err != nil {
		log.Fatal(err)
	}
}
