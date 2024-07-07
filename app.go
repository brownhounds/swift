package swift

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pb33f/libopenapi"
)

type Swift struct {
	middlewares []Middleware
	handlers    map[string]*HandlerValue
	serverMux   *http.ServeMux
	context     *Context
	server      http.Server
}

func New() *Swift {
	return &Swift{
		serverMux:   http.NewServeMux(),
		handlers:    make(map[string]*HandlerValue),
		middlewares: make([]Middleware, 0),
		context:     &ContextValue,
	}
}

func (r *Swift) Get(path string, handler Handler) *HandlerValue {
	p := BuildAndValidatePath(path)
	return r.MakeHandler(http.MethodGet, p, handler, nil)
}

func (r *Swift) Post(path string, handler Handler) *HandlerValue {
	p := BuildAndValidatePath(path)
	return r.MakeHandler(http.MethodPost, p, handler, nil)
}

func (r *Swift) Put(path string, handler Handler) *HandlerValue {
	p := BuildAndValidatePath(path)
	return r.MakeHandler(http.MethodPut, p, handler, nil)
}

func (r *Swift) Patch(path string, handler Handler) *HandlerValue {
	p := BuildAndValidatePath(path)
	return r.MakeHandler(http.MethodPatch, p, handler, nil)
}

func (r *Swift) Delete(path string, handler Handler) *HandlerValue {
	p := BuildAndValidatePath(path)
	return r.MakeHandler(http.MethodDelete, p, handler, nil)
}

func (r *Swift) Custom404(c func(w http.ResponseWriter)) {
	r.context.notFound = c
}

func (r *Swift) Custom500(c func(w http.ResponseWriter)) {
	r.context.internalServerError = c
}

func (r *Swift) Custom405(c func(w http.ResponseWriter)) {
	r.context.methodNotAllowed = c
}

func (r *Swift) Middleware(m ...Middleware) {
	r.middlewares = append(r.middlewares, m...)
}

func (r *Swift) OApiSchema(path string) {
	schema, err := os.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	document, err := libopenapi.NewDocument(schema)
	if err != nil {
		panic(err)
	}

	r.context.schema = &document
	BuiltInMiddlewares = append(BuiltInMiddlewares, ValidateOApiSchemaMiddleware)
}

func (r *Swift) Serve(host, port string) {
	Boot(r)

	r.server = http.Server{
		Addr: fmt.Sprintf("%s:%s", host, port),
		Handler: MakeMiddlewareStack(
			r.serverMux,
			Prepend(r.middlewares, BuiltInMiddlewares...)),
	}

	log.Printf("Listening on: %s:%s", host, port)

	if err := r.server.ListenAndServe(); err != nil {
		panic(err.Error())
	}
}
