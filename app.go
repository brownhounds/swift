package swift

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/pb33f/libopenapi"
)

type Swift struct {
	middlewares []Middleware
	handlers    map[string]*HandlerValue
	serverMux   *http.ServeMux
	context     *Context
	server      http.Server
}

type TLS struct {
	certificate string
	key         string
}

func New() *Swift {
	return &Swift{
		serverMux:   http.NewServeMux(),
		handlers:    make(map[string]*HandlerValue),
		middlewares: make([]Middleware, 0),
		context:     &ContextValue,
	}
}

func (r *Swift) Handle(pth string, handler http.Handler) {
	r.serverMux.Handle(pth, handler)
}

func (r *Swift) Get(pth string, handler http.HandlerFunc) *HandlerValue {
	p := BuildAndValidatePath(pth)
	return r.MakeHandler(http.MethodGet, p, handler, nil)
}

func (r *Swift) Post(pth string, handler http.HandlerFunc) *HandlerValue {
	p := BuildAndValidatePath(pth)
	return r.MakeHandler(http.MethodPost, p, handler, nil)
}

func (r *Swift) Put(pth string, handler http.HandlerFunc) *HandlerValue {
	p := BuildAndValidatePath(pth)
	return r.MakeHandler(http.MethodPut, p, handler, nil)
}

func (r *Swift) Patch(pth string, handler http.HandlerFunc) *HandlerValue {
	p := BuildAndValidatePath(pth)
	return r.MakeHandler(http.MethodPatch, p, handler, nil)
}

func (r *Swift) Delete(pth string, handler http.HandlerFunc) *HandlerValue {
	p := BuildAndValidatePath(pth)
	return r.MakeHandler(http.MethodDelete, p, handler, nil)
}

func (r *Swift) Custom404(c func(w http.ResponseWriter)) {
	r.context.notFound = c
}

func (r *Swift) Custom500(c func(w http.ResponseWriter)) {
	r.context.internalServerError = c
}

func (r *Swift) Middleware(m ...Middleware) {
	r.middlewares = append(r.middlewares, m...)
}

func (r *Swift) OApiValidator(pathToSchema string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	schema, err := os.ReadFile(path.Join(dir, pathToSchema))
	if err != nil {
		panic(err)
	}

	document, err := libopenapi.NewDocument(schema)
	if err != nil {
		panic(err)
	}

	r.context.schema = &document
	BuiltInMiddlewares = append(BuiltInMiddlewares, ValidateOApiSchemaMiddleware)
}

func (r *Swift) SetApiPrefix(prefix string) {
	p := BuildAndValidatePath(prefix)
	r.context.apiPrefix = p
}

func (r *Swift) ApiPrefix() *Group {
	return r.Group(r.context.apiPrefix)
}

func (r *Swift) AddTLS(crt, key string) {
	r.context.tls = &TLS{
		certificate: crt,
		key:         key,
	}
}

func (r *Swift) OnBoot(f func()) {
	r.context.onBoot = f
}

func (r *Swift) AddCorsMiddleware(f func(next http.Handler) http.Handler) {
	r.context.cors = f
}

func (r *Swift) Serve(host, port string) {
	Boot(r)

	if r.context.onBoot != nil {
		r.context.onBoot()
	}

	handler := MakeMiddlewareStack(r.serverMux, Prepend(r.middlewares, BuiltInMiddlewares...))
	if r.context.cors != nil {
		handler = r.context.cors(handler)
	}

	r.server = http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: handler,
	}

	log.Printf("Listening on: %s:%s", host, port)

	if r.context.tls != nil {
		if err := r.server.ListenAndServeTLS(r.context.tls.certificate, r.context.tls.key); err != nil {
			panic(err)
		}
	} else {
		if err := r.server.ListenAndServe(); err != nil {
			panic(err)
		}
	}
}
