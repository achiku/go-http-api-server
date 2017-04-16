package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func testNewApp(t *testing.T) *App {
	var logger *log.Logger
	if testing.Verbose() {
		logger = log.New(os.Stdout, "[test log] ", log.LstdFlags)
	} else {
		logger = log.New(ioutil.Discard, "[null log] ", log.LstdFlags)
	}
	return &App{
		Name:   "my-test-server",
		Host:   "test-host",
		Logger: logger,
	}
}

func TestGreeting(t *testing.T) {
	app := testNewApp(t)
	req := httptest.NewRequest(http.MethodGet, "/api/hello", nil)
	w := httptest.NewRecorder()
	status, res, err := app.Greeting(w, req)
	if err != nil {
		t.Fatal(err)
	}
	if status != http.StatusOK {
		t.Fatalf("want %d got %d", http.StatusOK, status)
	}
	gt, ok := res.(*Greeting)
	if !ok {
		t.Fatalf("want type Greeting got %s", reflect.TypeOf(res))
	}
	if expected := "anonymous"; gt.Name != expected {
		t.Errorf("want %s got %s", expected, gt.Name)
	}
	if expected := "hello!"; gt.Message != expected {
		t.Errorf("want %s got %s", expected, gt.Message)
	}
}
