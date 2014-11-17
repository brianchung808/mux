package mux

import (
	"net/http"
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

type Route struct {
	path      string
	endpoints verbHandlerMap
}

type Router struct {
	// (path_URI -> route) map.
	routes map[string]*Route
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]*Route),
	}
}

func (r *Router) Handle(path string, verb string, handler http.Handler) {
	route := r.routes[path]

	if route == nil {
		route := &Route{
			path:      path,
			endpoints: make(verbHandlerMap),
		}
		// set the new route
		r.routes[path] = route
	}

	// set the handler
	route.endpoints[verb] = handler
}

// satisfy Handler interface
// handles all requests & delegate to other routes.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	// find the corresponding Route in Router & call it's handler.
	route := r.routes[path]

	// the handler we will be delegating to
	var handler http.Handler

	if route == nil {
		// route not found
		handler = http.NotFoundHandler()
	}

	if handler = route.endpoints[method]; handler == nil {
		// handler not found
		handler = http.NotFoundHandler()
	}

	handler.ServeHTTP(w, req)
}
