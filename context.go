package swift

import (
	"net/http"

	"github.com/pb33f/libopenapi"
)

type Context struct {
	notFound            func(w http.ResponseWriter)
	internalServerError func(w http.ResponseWriter)
	cors                func(next http.Handler) http.Handler
	schema              *libopenapi.Document
	tls                 *TLS
	swagger             *SwaggerServer
	staticServer        *StaticServer
	onBoot              func()
	apiPrefix           string
}

var ContextValue = Context{
	apiPrefix:           "",
	notFound:            nil,
	internalServerError: nil,
	schema:              nil,
	tls:                 nil,
	swagger:             nil,
	staticServer:        nil,
	onBoot:              nil,
	cors:                nil,
}

func GetContext() *Context {
	return &ContextValue
}
