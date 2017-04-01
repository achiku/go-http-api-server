package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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
	app.Greeting(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("want %d got %d", http.StatusOK, w.Code)
	}
	var res Greeting
	decoder := json.NewDecoder(w.Body)
	if err := decoder.Decode(&res); err != nil {
		t.Fatal(err)
	}
	if expected := "anonymous"; res.Name != expected {
		t.Errorf("want %s got %s", expected, res.Name)
	}
	if expected := "hello!"; res.Message != expected {
		t.Errorf("want %s got %s", expected, res.Message)
	}
}
