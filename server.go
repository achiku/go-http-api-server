package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/pkg/errors"
	"github.com/rs/xlog"
)

// App application
type App struct {
	Name string
}

// InternalHandler internal
type InternalHandler struct {
	h func(w http.ResponseWriter, r *http.Request) (int, interface{}, error)
}

// ServeHTTP serve
func (ih InternalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := xlog.FromRequest(r)
	encoder := json.NewEncoder(w)
	reqInfo := xlog.F{"http_request": r}

	statusCode, res, err := ih.h(w, r)
	if err != nil {
		logger.Error(err, reqInfo)
		w.WriteHeader(statusCode)
		encoder.Encode(res)
		return
	}
	w.WriteHeader(statusCode)
	encoder.Encode(res)
	return
}

func (app *App) hello(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	logger := xlog.FromRequest(r)
	logger.Info("hello handler")
	res, err := HelloService(r.Context(), "")
	if err != nil {
		return http.StatusInternalServerError, nil, errors.Wrap(err, "HelloService failed")
	}
	return http.StatusOK, res, nil
}

func (app *App) helloWithName(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	logger := xlog.FromRequest(r)
	logger.Info("hello wth name handler")
	val := mux.Vars(r)
	res, err := HelloService(r.Context(), val["name"])
	if err != nil {
		return http.StatusInternalServerError, nil, errors.Wrap(err, "HelloService failed")
	}
	return http.StatusOK, res, nil
}

func main() {
	host, _ := os.Hostname()
	logConf := xlog.Config{
		Fields: xlog.F{
			"role": "my-service",
			"host": host,
		},
		Output: xlog.NewOutputChannel(xlog.NewConsoleOutput()),
	}

	c := alice.New(
		xlog.NewHandler(logConf),
		xlog.MethodHandler("method"),
		xlog.URLHandler("url"),
		xlog.UserAgentHandler("user_agent"),
		xlog.RefererHandler("referer"),
		xlog.RequestIDHandler("req_id", "Request-Id"),
		accessLoggingMiddleware,
	)
	app := App{Name: "my-service"}

	// for gorilla/mux
	router := mux.NewRouter()
	r := router.PathPrefix("/api").Subrouter()
	r.Methods("GET").Path("/hello").Handler(c.Then(InternalHandler{h: app.hello}))
	r.Methods("GET").Path("/hello/staticName").Handler(c.Then(InternalHandler{h: app.hello}))
	r.Methods("GET").Path("/hello/{name}").Handler(c.Then(InternalHandler{h: app.helloWithName}))

	xlog.Info("start server")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
