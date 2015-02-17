mux
===

HTTP routing library compatible with `http.Handler` interface.

Testing out go-lang.

Usage:

```go

// get a new router 
r := mux.NewRouter()

// register route with function literal
r.HandleFunc("/restaurant", "GET", func(w ResponseWriter, r *Request) {
	w.Write([]byte(`GET`))
})

// register multiple routes specified in Endpoint struct
r.HandleAll("/restaurant', Endpoint{
	Get: func(w http.ResponseWriter, req * http.Request) {},
	Post: func(w http.ResponseWriter, req * http.Request) {},
	...
})

type Hello string
// implement http.Handler for Hello
func (h Hello) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(h))
}

// register route with type that implements http.Handler
r.Handle("/restaurant", Hello("World"))

```
