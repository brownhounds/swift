package swift

import (
	"net/http"

	"github.com/pb33f/libopenapi"
)

type Context struct {
	notFound            func(w http.ResponseWriter)
	methodNotAllowed    func(w http.ResponseWriter)
	internalServerError func(w http.ResponseWriter)
	schema              *libopenapi.Document
	tls                 *TLS
}

var ContextValue = Context{
	notFound:            nil,
	methodNotAllowed:    nil,
	internalServerError: nil,
	schema:              nil,
	tls:                 nil,
}

func GetContext() *Context {
	return &ContextValue
}
