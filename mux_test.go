package mux

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
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

	router.HandleFunc("/test/", GET, func(w http.ResponseWriter, req *http.Request) {})

	route := router.routes["/test/"]

	if assert.NotNil(t, route, "Route is missing") {
		assert.Equal(t, "/test/", route.path, "Path incorrect")
	}

	endpoints := route.endpoints

	if assert.NotNil(t, endpoints, "Endpoints are missing") {
		assert.NotNil(t, endpoints[GET], "Endpoint GET missing")
	}
}

func TestRouteResponse(t *testing.T) {
	router := NewRouter()

	// test recorder that implements http.ResponseWriter
	w := httptest.NewRecorder()

	router.HandleFunc("/test/", GET, func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`yolo`))
	})

	req, err := http.NewRequest("GET", "/test/", nil)

	if err != nil {
		t.Fail()
	}

	// call the router's ServeHTTP directly
	router.ServeHTTP(w, req)

	assert.Equal(t, "yolo", w.Body.String(), "Incorrect Body response")
}

func TestNonExistingRouteResponse(t *testing.T) {
	router := NewRouter()

	verbs := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

	for _, verb := range verbs {
		req, err := http.NewRequest(verb, "/test/", nil)
		if err != nil {
			t.Fail()
		}

		expectNotFoundHandler(t, router, req)
	}
}

func TestNonExistingMethodResponse(t *testing.T) {
	router := NewRouter()

	router.HandleFunc("/test1/", GET, func(w http.ResponseWriter, req *http.Request) {})

	req, err := http.NewRequest("POST", "/test1/", nil)

	if err != nil {
		t.Fail()
	}

	expectNotFoundHandler(t, router, req)
}

func TestMultipleRouteData(t *testing.T) {
	router := NewRouter()

	router.HandleFunc("/test1/", GET, func(w http.ResponseWriter, req *http.Request) {})

	router.HandleFunc("/test2/", GET, func(w http.ResponseWriter, req *http.Request) {})

	expected := []string{"/test1/", "/test2/"}

	testRoutePathInfo(t, expected, router)

	validVerbs := []int{GET}
	testVerbs(t, validVerbs, router)
}

func TestMultipleVerbData(t *testing.T) {
	router := NewRouter()

	validVerbs := []int{GET, POST, PATCH}
	for _, verb := range validVerbs {
		router.HandleFunc("/test1/", verb, func(w http.ResponseWriter, req *http.Request) {})
	}

	testVerbs(t, validVerbs, router)

	expected := []string{"/test1/"}
	testRoutePathInfo(t, expected, router)

}

func TestCleanupPath(t *testing.T) {
	paths := []string{
		"/hello/hi",
		"/hello",
		"/hello/hi      ",
		"/hello     ",
		"",
		" ",
	}

	expPaths := []string{
		"/hello/hi/",
		"/hello/",
		"/hello/hi/",
		"/hello/",
		"/",
		"/",
	}

	for i, path := range paths {
		newPath := cleanupPath(strings.NewReader(path))
		exp := expPaths[i]

		assert.Equal(t, exp, newPath, "Paths not equal")
	}
}

func TestHandleAll(t *testing.T) {
	router := NewRouter()

	router.HandleAll("/test/", Endpoint{
		Get:  func(w http.ResponseWriter, req *http.Request) { w.Write([]byte(`GET`)) },
		Post: func(w http.ResponseWriter, req *http.Request) { w.Write([]byte(`POST`)) },
	})

	validVerbs := []int{GET, POST}

	// test if routes registered
	testVerbs(t, validVerbs, router)

	expected := []string{"/test/"}
	testRoutePathInfo(t, expected, router)

	// test recorder that implements http.ResponseWriter
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/test/", nil)

	if err != nil {
		t.Fail()
	}

	// call the router's ServeHTTP directly
	router.ServeHTTP(w, req)

	assert.Equal(t, "GET", w.Body.String(), "Incorrect Body response")

	w = httptest.NewRecorder()

	req, err = http.NewRequest("POST", "/test/", nil)

	if err != nil {
		t.Fail()
	}

	// call the router's ServeHTTP directly
	router.ServeHTTP(w, req)

	assert.Equal(t, "POST", w.Body.String(), "Incorrect Body response")
}

//************
// Helpers
//************

func testRoutePathInfo(t *testing.T, expected []string, router *router) {
	for _, exp := range expected {
		route, ok := router.routes[exp]

		if assert.True(t, ok, "Route missing") {
			assert.Equal(t, exp, route.path, "Path missing")
		}
	}
}

func testVerbs(t *testing.T, validVerbs []int, router *router) {
	for path, route := range router.routes {
		assert.Equal(t, route.path, path, "Path key not equal to route.path it is pointing to")

		endpoints := route.endpoints

		if assert.NotNil(t, endpoints, "Endpoints missing") {
			for _, verb := range validVerbs {
				assert.NotNil(t, endpoints[verb], "Endpoint missing")
			}
		}
	}
}

// test if handler is correctly Not Found Handler
func expectNotFoundHandler(t *testing.T, router *router, req *http.Request) {
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req)

	w2 := httptest.NewRecorder()
	http.NotFound(w2, req)

	assert.Equal(t, w2.Body.String(), w1.Body.String(), "Incorrect body response")
}
