package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/xlog"
)

// App application
type App struct {
	Name string
}

// InternalHandler internal
type InternalHandler struct {
	h func(w http.ResponseWriter, r *http.Request)
}

// ServeHTTP serve
func (ih InternalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ih.h(w, r)
	return
}

func (app *App) hello(w http.ResponseWriter, r *http.Request) {
	logger := xlog.FromRequest(r)
	logger.Info("hello handler")
	log.Println("this is usual logger")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello")
	return
}

func (app *App) helloWithName(w http.ResponseWriter, r *http.Request) {
	logger := xlog.FromRequest(r)
	logger.Info("hello handler")
	log.Println("this is usual logger")
	val := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello, %s", val["name"])
	return
}

func f(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(h.ServeHTTP)
}

func main() {
	host, _ := os.Hostname()
	conf := xlog.Config{
		Fields: xlog.F{
			"role": "my-service",
			"host": host,
		},
		Output: xlog.NewOutputChannel(xlog.NewConsoleOutput()),
	}

	c := alice.New(
		xlog.NewHandler(conf),
		xlog.MethodHandler("method"),
		xlog.URLHandler("url"),
		xlog.UserAgentHandler("user_agent"),
		xlog.RefererHandler("referer"),
		xlog.RequestIDHandler("req_id", "Request-Id"),
		accessLoggingMiddleware,
	)
	app := App{Name: "my-service"}
	// for httptreemux
	// router := httptreemux.New()
	// r := router.NewGroup("/api").UsingContext()
	// r.GET("/hello", c.Then(InternalHandler{h: app.hello}))
	// r.GET("/hello", http.HandlerFunc(c.Then(InternalHandler{h: app.hello})))
	// h := c.Then(InternalHandler{h: app.hello})
	// r.GET("/hello", http.HandlerFunc(h.ServeHTTP))
	// r.GET("/hello", f(c.Then(InternalHandler{h: app.hello})))
	// r.GET("/hello", c.Then(InternalHandler{h: app.hello}).ServeHTTP)

	// for gorilla/mux
	router := mux.NewRouter()
	r := router.PathPrefix("/api").Subrouter()
	r.Methods("GET").Path("/hello").Handler(c.Then(InternalHandler{h: app.hello}))
	r.Methods("GET").Path("/hello/staticName").Handler(c.Then(InternalHandler{h: app.hello}))
	r.Methods("GET").Path("/hello/{name}").Handler(c.Then(InternalHandler{h: app.helloWithName}))

	xlog.Info("xlog")
	xlog.Infof("chain: %+v", c)
	log.Println("start server")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
