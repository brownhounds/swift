package swift

import (
	"net/http"

	"github.com/brownhounds/swift/res"
)

type HandlerValue struct {
	method      string
	path        string
	handler     http.HandlerFunc
	group       *Group
	middlewares []Middleware
}

func (h *HandlerValue) Middleware(m ...Middleware) {
	h.middlewares = append(h.middlewares, m...)
}

func (r *Swift) MakeHandler(m, path string, handler http.HandlerFunc, group *Group) *HandlerValue {
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
