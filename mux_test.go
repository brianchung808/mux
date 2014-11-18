package mux

import (
	"github.com/stretchr/testify/assert"
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

	assert.NotNil(t, router, "For whatever reason, router is nil")
	assert.NotNil(t, router.routes, "[]Routes not correctly initialized")
}

func TestRoute(t *testing.T) {
	router := NewRouter()

	router.HandleFunc("/test", "GET", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`yolo`))
	})

	route := router.routes["/test"]

	if assert.NotNil(t, route, "Route is missing") {
		assert.Equal(t, route.path, "/test", "Path incorrect")
	}

	endpoints := route.endpoints

	if assert.NotNil(t, endpoints, "Endpoints are missing") {
		assert.NotNil(t, endpoints["GET"], "Endpoint GET missing")
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

	if assert.NotNil(t, route, "Route missing") {
		assert.Equal(t, route.path, "/test1", "Path missing")
	}

	endpoints := route.endpoints

	if assert.NotNil(t, endpoints, "Endpoints missing") {
		assert.NotNil(t, endpoints["GET"], "Endpoint GET missing")
		assert.Nil(t, endpoints["POST"], "Unregistered endpoint not nil")
	}

	route = router.routes["/test2"]

	route = router.routes["/test2"]

	if assert.NotNil(t, route, "Route missing") {
		assert.NotEqual(t, route.path, "/test1", "Path missing")
	}

	endpoints = route.endpoints

	if assert.NotNil(t, endpoints, "Endpoints missing") {
		assert.NotNil(t, endpoints["GET"], "Endpoint GET missing")
		assert.Nil(t, endpoints["POST"], "Unregistered endpoint not nil")
	}

}
