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

// verb -> handler
type verbHandlerMap map[string]http.Handler

type route struct {
	path      string
	endpoints verbHandlerMap
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
	r.handle(path, "GET", handler)
}

func (r *router) Post(path string, handler http.Handler) {
	r.handle(path, "POST", handler)
}

func (r *router) Delete(path string, handler http.Handler) {
	r.handle(path, "DELETE", handler)
}

func (r *router) Put(path string, handler http.Handler) {
	r.handle(path, "PUT", handler)
}

func (r *router) Patch(path string, handler http.Handler) {
	r.handle(path, "PATCH", handler)
}

func (r *router) handle(path string, verb string, handler http.Handler) {
	// clean up path
	path = cleanupPath(strings.NewReader(path))
	currentRoute := r.routes[path]

	if currentRoute == nil {
		currentRoute = &route{
			path:      path,
			endpoints: make(verbHandlerMap),
		}
		// set the new route
		r.routes[path] = currentRoute
	}

	// set the handler
	currentRoute.endpoints[verb] = handler
}

type handlerFunc func(http.ResponseWriter, *http.Request)

// to wrap handlerFuncs
type wrapper handlerFunc

// wrapper implements http.Handler interface & delegates to its handler
func (w wrapper) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	w(wri, req)
}

func (r *router) HandleFunc(path string, verb string, handler handlerFunc) {
	r.handle(path, verb, wrapper(handler))
}

// satisfy Handler interface
// handles all requests & delegate to other routes.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	// find the corresponding Route in Router & call it's handler.
	route := r.routes[path]

	// the handler we will be delegating to
	var handler http.Handler

	if route == nil {
		// route not found
		handler = http.NotFoundHandler()

	} else {
		// route exists
		endpoints := route.endpoints
		if endpoints != nil {
			var ok bool
			if handler, ok = endpoints[method]; !ok {
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
