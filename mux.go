package mux

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

/*
Usage:
r := mux.NewRouter()
r.Handle("/restaurant", "GET", func(w ResponseWriter, r *Request) {

})

r.HandleAll("/restaurant', Endpoint{
	GET: Handler,
	POST: Handler,
	...
})


r.Handle("/restaurant")

*/

// enum for common http verbs
const (
	GET = iota
	POST
	PUT
	PATCH
	DELETE
	OPTIONS
	HEAD
	NUM_VERBS
)

// verb -> handler
type verbHandlerList []http.Handler

type route struct {
	path      string
	endpoints verbHandlerList
}

type router struct {
	// (path_URI -> route) map.
	routes map[string]*route
}

func NewRouter() *router {
	return &router{
		routes: make(map[string]*route),
	}
}

func (r *router) Get(path string, handler http.Handler) {
	r.handle(path, GET, handler)
}

func (r *router) Post(path string, handler http.Handler) {
	r.handle(path, POST, handler)
}

func (r *router) Delete(path string, handler http.Handler) {
	r.handle(path, DELETE, handler)
}

func (r *router) Put(path string, handler http.Handler) {
	r.handle(path, PUT, handler)
}

func (r *router) Patch(path string, handler http.Handler) {
	r.handle(path, PATCH, handler)
}

func (r *router) handle(path string, verb int, handler http.Handler) {
	// clean up path
	path = cleanupPath(strings.NewReader(path))
	currentRoute := r.routes[path]

	// check if route exists or not
	if currentRoute == nil {
		currentRoute = &route{
			path:      path,
			endpoints: make(verbHandlerList, NUM_VERBS, NUM_VERBS),
		}
		// set the new route
		r.routes[path] = currentRoute
	}

	// set the handler
	currentRoute.endpoints[verb] = handler
}

// for function literals being passed as route handler
type routeHandlerFunc func(http.ResponseWriter, *http.Request)

// handlerFunc implements http.Handler interface & delegates to its handler
func (r routeHandlerFunc) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	r(wri, req)
}

func (r *router) HandleFunc(path string, verb int, handler routeHandlerFunc) {
	r.handle(path, verb, routeHandlerFunc(handler))
}

// satisfy http.Handler interface
// router handles all requests & delegate to other routes.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
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
		endpoints := route.endpoints
		if endpoints != nil {
			if handler = endpoints[methodEnum]; handler == nil {
				// handler not found
				handler = http.NotFoundHandler()
			}
		} else {
			handler = http.NotFoundHandler()
		}
	}

	handler.ServeHTTP(w, req)
}

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
