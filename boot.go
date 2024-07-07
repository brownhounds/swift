package swift

import (
	"net/http"
)

func Boot(r *Swift) {
	initializeNotFoundHandler(r)
	initializeHandlers(r)
	initializeMethodNotAllowedHandlers(r)
}

func initializeNotFoundHandler(r *Swift) {
	r.serverMux.HandleFunc("/", NotFoundHandler)
}

func initializeHandlers(r *Swift) {
	for _, v := range r.handlers {
		stack := MakeMiddlewareStack(http.HandlerFunc(v.handler), v.middlewares)
		if v.group != nil {
			groupStack := MakeMiddlewareStack(stack, v.group.middlewares)
			r.serverMux.Handle(v.method+" "+v.path, groupStack)
		} else {
			r.serverMux.Handle(v.method+" "+v.path, stack)
		}
	}
}

func initializeMethodNotAllowedHandlers(r *Swift) {
	group := make(map[string]struct{})
	for _, v := range r.handlers {
		if v.path != "/" {
			group[v.path] = struct{}{}
		}
	}

	for path := range group {
		r.serverMux.HandleFunc(path, MethodNotAllowedHandler)
	}
}
