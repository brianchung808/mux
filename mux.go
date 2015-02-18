// Package mux provides a http request multiplexer
package mux

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

// enum for common http methods
const (
	GET = iota
	POST
	PUT
	PATCH
	DELETE
	OPTIONS
	HEAD
	NUM_METHODS
)

// struct for user's to register all of their endpoints at once
type Handlers struct {
	Get     http.HandlerFunc
	Post    http.HandlerFunc
	Put     http.HandlerFunc
	Patch   http.HandlerFunc
	Delete  http.HandlerFunc
	Options http.HandlerFunc
	Head    http.HandlerFunc
}

// type router holds the routes that are being handled
type router struct {
	// (path_URI -> route) map.
	routes map[string]*route
}

// return a new router
func NewRouter() *router {
	return &router{
		routes: make(map[string]*route),
	}
}

// register a handler for the specified path for the method
func (r router) handle(path string, method int, handler http.Handler) {
	// clean up path
	path = cleanupPath(strings.NewReader(path))
	currentRoute := r.routes[path]

	// check if route exists or not
	if currentRoute == nil {
		currentRoute = &route{
			path:      path,
			endpoints: make([]http.Handler, NUM_METHODS, NUM_METHODS),
		}
		// set the new route
		r.routes[path] = currentRoute
	}

	// set the handler
	currentRoute.endpoints[method] = handler
}

// register handler for Get
func (r router) Get(path string, handler http.Handler) {
	r.handle(path, GET, handler)
}

// register handler for Post
func (r router) Post(path string, handler http.Handler) {
	r.handle(path, POST, handler)
}

// register handler for Delete
func (r router) Delete(path string, handler http.Handler) {
	r.handle(path, DELETE, handler)
}

// register handler for Put
func (r router) Put(path string, handler http.Handler) {
	r.handle(path, PUT, handler)
}

// register handler for Patch
func (r router) Patch(path string, handler http.Handler) {
	r.handle(path, PATCH, handler)
}

// register handler for Options
func (r router) Options(path string, handler http.Handler) {
	r.handle(path, OPTIONS, handler)
}

// register handler for Head
func (r router) Head(path string, handler http.Handler) {
	r.handle(path, HEAD, handler)
}

// register handler for specified method
func (r router) HandleFunc(path string, method int, handler http.HandlerFunc) {
	r.handle(path, method, handler)
}

// register multiple handlers at once with an Handlers struct
func (r router) HandleAll(path string, h Handlers) {
	r.Get(path, h.Get)
	r.Post(path, h.Post)
	r.Put(path, h.Put)
	r.Patch(path, h.Patch)
	r.Delete(path, h.Delete)
	r.Options(path, h.Options)
	r.Head(path, h.Head)
}

// satisfy http.Handler interface
// router handles all requests & delegate to other routes.
func (r router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := cleanupPath(strings.NewReader(req.URL.Path))
	method := req.Method

	var methodEnum int

	switch method {
	case "GET":
		methodEnum = GET
	case "POST":
		methodEnum = POST
	case "PUT":
		methodEnum = PUT
	case "DELETE":
		methodEnum = DELETE
	case "PATCH":
		methodEnum = PATCH
	}

	// find the corresponding Route in Router & call it's handler.
	route, ok := r.routes[path]

	// the handler we will be delegating to
	var handler http.Handler

	if !ok {
		// route not found
		handler = http.NotFoundHandler()
	} else {
		// route exists
		// delegate to route.handle

		handler = route.handler(methodEnum)
	}

	handler.ServeHTTP(w, req)
}

// helper to cleanup user inputted path string
func cleanupPath(path io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(path)

	pathBytes := buf.Bytes()

	// trim
	pathBytes = bytes.Trim(pathBytes, " ")

	len := len(pathBytes)
	// empty case
	if len == 0 {
		return "/"
	}

	// check for trailing '/'
	if pathBytes[len-1] != '/' {
		pathBytes = append(pathBytes, '/')
	}

	return string(pathBytes)
}
