package mux

import (
	"net/http"
)

// type route represents a route to the path with the specified handlers
type route struct {
	path      string
	endpoints []http.Handler
}

// return the correct handler for the given method for the route
func (r route) handler(method int) (handler http.Handler) {
	endpoints := r.endpoints
	if endpoints != nil {
		if handler = endpoints[method]; handler == nil {
			// handler not found
			handler = http.NotFoundHandler()
		}
	} else {
		handler = http.NotFoundHandler()
	}

	return
}
