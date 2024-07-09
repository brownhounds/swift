package swift

import (
	"net/http"
)

func Boot(r *Swift) {
	initializeHealthCheck(r)
	initializeSwaggerServer(r)
	initializeNotFoundHandler(r)
	initializeHandlers(r)
	initializeMethodNotAllowedHandlers(r)
}

func initializeHealthCheck(r *Swift) {
	r.serverMux.HandleFunc("GET /health", HealthCheckHandler)
}

func initializeSwaggerServer(r *Swift) {
	if r.context.swagger != nil && r.context.swagger.serve {
		path := r.context.swagger.path
		fileServer := http.FileServer(http.Dir("./" + r.context.swagger.staticDir + "/"))
		r.serverMux.Handle(path, http.StripPrefix(path, fileServer))
	}
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
