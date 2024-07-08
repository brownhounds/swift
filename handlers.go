package swift

import (
	"net/http"

	res "github.com/brownhounds/swift/response"
)

type Handler func(http.ResponseWriter, *http.Request)

type HandlerValue struct {
	method      string
	path        string
	handler     Handler
	group       *Group
	middlewares []Middleware
}

func (h *HandlerValue) Middleware(m ...Middleware) {
	h.middlewares = append(h.middlewares, m...)
}

func (r *Swift) MakeHandler(m, path string, handler Handler, group *Group) *HandlerValue {
	value := &HandlerValue{
		method:      m,
		path:        path,
		handler:     handler,
		middlewares: make([]Middleware, 0),
		group:       group,
	}
	r.handlers[m+" "+path] = value
	return value
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext()

	if ctx.methodNotAllowed != nil {
		ctx.methodNotAllowed(w)
	} else {
		res.ApiError(w, http.StatusMethodNotAllowed)
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext()

	if ctx.notFound != nil {
		ctx.notFound(w)
	} else {
		res.ApiError(w, http.StatusNotFound)
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	res.Json(w, http.StatusOK, res.Map{
		"status": "OK",
	})
}
