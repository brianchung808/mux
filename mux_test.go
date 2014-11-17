package mux

import (
	"net/http"
	"testing"
)

type test struct {
	msg string
}

func (t *test) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(t.msg))
}

// var router *Router

// func init() {
// 	router = NewRouter()
// }

func TestRouter(t *testing.T) {
	router := NewRouter()

	if router.routes == nil {
		t.Error("Failed to init Router")
	}
}

func TestRoute(t *testing.T) {
	router := NewRouter()

	router.HandleFunc("/test", "GET", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`yolo`))
	})

	route := router.routes["/test"]

	if route == nil {
		t.Error("Route missing")
	}

	if route.path != "/test" {
		t.Error("Path missing")
	}

	endpoints := route.endpoints

	if endpoints == nil {
		t.Error("Endpoints missing")
	}

	if endp := endpoints["GET"]; endp == nil {
		t.Error("Endpoint GET missing")
	}
}

func TestMultipleRoutes(t *testing.T) {
	router := NewRouter()

	router.HandleFunc("/test1", "GET", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`yolo`))
	})

	router.HandleFunc("/test2", "GET", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`yolo`))
	})

	route := router.routes["/test1"]

	if route == nil {
		t.Error("Route missing")
	}

	if route.path != "/test1" {
		t.Error("Path missing")
	}

	endpoints := route.endpoints

	if endpoints == nil {
		t.Error("Endpoints missing")
	}

	if endp := endpoints["GET"]; endp == nil {
		t.Error("Endpoint GET missing")
	}

	route = router.routes["/test2"]

	if route == nil {
		t.Error("Route missing")
	}

	if route.path != "/test2" {
		t.Error("Path missing")
	}

	endpoints = route.endpoints

	if endpoints == nil {
		t.Error("Endpoints missing")
	}

	if endp := endpoints["GET"]; endp == nil {
		t.Error("Endpoint GET missing")
	}
}
