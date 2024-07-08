package swift

import (
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	res "github.com/brownhounds/swift/response"

	validator "github.com/pb33f/libopenapi-validator"
)

type Middleware func(http.Handler) http.Handler

func MiddlewareStack(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

var BuiltInMiddlewares = []Middleware{RecoverMiddleware}

func MakeMiddlewareStack(h http.Handler, m []Middleware) http.Handler {
	if len(m) > 0 {
		return MiddlewareStack(h, m...)
	}

	return h
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v\n", err)
				log.Println(string(debug.Stack()))

				ctx := GetContext()

				if ctx.internalServerError != nil {
					ctx.internalServerError(w)
				} else {
					ApiError(w, http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func ValidateOApiSchemaMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := GetContext()

		if ctx.swagger != nil && strings.HasPrefix(r.URL.Path, ctx.swagger.path) {
			next.ServeHTTP(w, r)
			return
		}

		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		highLevelValidator, validatorErrs := validator.NewValidator(*ctx.schema)
		if len(validatorErrs) > 0 {
			panic(validatorErrs[0])
		}

		requestValid, validationErrors := highLevelValidator.ValidateHttpRequestSync(r)

		if !requestValid {
			for _, v := range validationErrors {
				if strings.HasSuffix(v.Message, "not found") {
					ApiError(w, http.StatusNotFound)
					break
				}

				res.Json(w, http.StatusBadRequest, res.Map{
					"status":  http.StatusBadRequest,
					"message": http.StatusText(http.StatusBadRequest),
					"reason":  v.Reason,
				})
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
