package mux

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type test struct {
	msg string
}

func (t *test) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(t.msg))
}

func TestRouter(t *testing.T) {
	router := NewRouter()

	assert.NotNil(t, router, "For whatever reason, router is nil")
	assert.NotNil(t, router.routes, "[]Routes not correctly initialized")
}

func TestRouteData(t *testing.T) {
	router := NewRouter()

	router.HandleFunc("/test", "GET", func(w http.ResponseWriter, req *http.Request) {})

	route := router.routes["/test"]

	if assert.NotNil(t, route, "Route is missing") {
		assert.Equal(t, "/test", route.path, "Path incorrect")
	}

	endpoints := route.endpoints

	if assert.NotNil(t, endpoints, "Endpoints are missing") {
		assert.NotNil(t, endpoints["GET"], "Endpoint GET missing")
	}
}

func TestRouteResponse(t *testing.T) {
	router := NewRouter()

	// test recorder that implements http.ResponseWriter
	w := httptest.NewRecorder()

	handler := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`yolo`))
	}

	router.HandleFunc("/test", "GET", handler)

	req, err := http.NewRequest("GET", "/test", nil)

	if err != nil {
		t.Fail()
	}

	handler(w, req)

	assert.Equal(t, "yolo", w.Body.String(), "Incorrect Body response")

}

func TestMultipleRouteData(t *testing.T) {
	router := NewRouter()

	router.HandleFunc("/test1", "GET", func(w http.ResponseWriter, req *http.Request) {})

	router.HandleFunc("/test2", "GET", func(w http.ResponseWriter, req *http.Request) {})

	expected := []string{"/test1", "/test2"}

	testRoutePathInfo(t, expected, router)

	validVerbs := []string{"GET"}
	testVerbs(t, validVerbs, router)
}

func TestMultipleVerbData(t *testing.T) {
	router := NewRouter()

	validVerbs := []string{"GET", "POST", "PATCH"}
	for _, verb := range validVerbs {
		router.HandleFunc("/test1", verb, func(w http.ResponseWriter, req *http.Request) {})
	}

	testVerbs(t, validVerbs, router)

	expected := []string{"/test1"}
	testRoutePathInfo(t, expected, router)
}

//************
// Helpers
//************

func testRoutePathInfo(t *testing.T, expected []string, router *Router) {
	for _, exp := range expected {
		route, ok := router.routes[exp]

		if assert.True(t, ok, "Route missing") {
			assert.Equal(t, exp, route.path, "Path missing")
		}
	}
}

func testVerbs(t *testing.T, validVerbs []string, router *Router) {
	for path, route := range router.routes {
		assert.Equal(t, route.path, path, "Path key not equal to route.path it is pointing to")

		endpoints := route.endpoints

		if assert.NotNil(t, endpoints, "Endpoints missing") {
			for _, verb := range validVerbs {
				assert.NotNil(t, endpoints[verb], "Endpoint "+verb+" missing")
			}
		}
	}
}
