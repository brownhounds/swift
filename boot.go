package swift

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/brownhounds/swift/res"
)

const (
	DEV_FRONTEND_PROXY = "DEV_FRONTEND_PROXY"
)

func Boot(r *Swift) {
	initializeHealthCheck(r)
	initializeRootStaticServer(r)
	initializeSwaggerServer(r)
	initializeNotFoundHandler(r)
	initializeHandlers(r)
}

func initializeHealthCheck(r *Swift) {
	r.serverMux.HandleFunc("GET /health", HealthCheckHandler)
}

func initializeSwaggerServer(r *Swift) {
	if r.context.swagger != nil && r.context.swagger.serve {
		path := r.context.swagger.path
		fileServer := http.FileServer(http.Dir("./" + r.context.swagger.staticDir + "/"))
		r.serverMux.Handle("GET "+path, http.StripPrefix(path, fileServer))
	}
}

func initializeRootStaticServer(r *Swift) {
	_, defined := os.LookupEnv(DEV_FRONTEND_PROXY)
	if defined {
		remote, err := url.Parse(os.Getenv(DEV_FRONTEND_PROXY))
		if err != nil {
			panic(err)
		}
		handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				log.Printf("DEV Proxy: %s", r.URL)
				r.Host = remote.Host
				p.ServeHTTP(w, r)
			}
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		r.serverMux.HandleFunc("GET /", handler(proxy))
		return
	}

	if r.context.staticServer != nil {
		path := r.context.staticServer.path
		staticDir := "./" + r.context.staticServer.staticDir + "/"
		fileServer := http.FileServer(http.Dir(staticDir))
		r.serverMux.Handle("GET "+path, r.blessedHandler(path, staticDir, http.StripPrefix(path, fileServer)))
	}
}

func initializeNotFoundHandler(r *Swift) {
	r.serverMux.HandleFunc("/", NotFoundHandler)
}

func initializeHandlers(r *Swift) {
	for _, v := range r.handlers {
		stack := MakeMiddlewareStack(
			blessedHandlerFunc(v.path, v.handler), v.middlewares,
		)
		if v.group != nil {
			groupStack := MakeMiddlewareStack(stack, v.group.middlewares)
			r.serverMux.Handle(v.method+" "+v.path, groupStack)
		} else {
			r.serverMux.Handle(v.method+" "+v.path, stack)
		}
	}
}

func (s *Swift) blessedHandler(path, staticDir string, handler http.Handler) http.Handler {
	if path != "/" {
		return handler
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join(staticDir, r.URL.Path)
		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
			if s.context.apiPrefix != "" && strings.HasPrefix(r.URL.Path, s.context.apiPrefix) {
				NotFoundHandler(w, r)
			} else {
				errorPage := filepath.Join(staticDir, "404.html")
				if _, err := os.Stat(errorPage); errors.Is(err, os.ErrNotExist) {
					NotFoundHandler(w, r)
					return
				}
				if s.context.staticServer.spa {
					res.HtmlTemplate(w, http.StatusOK, filepath.Join(staticDir, "index.html"), nil)
					return
				}
				res.HtmlTemplate(w, http.StatusNotFound, errorPage, nil)
			}
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func blessedHandlerFunc(path string, handler http.HandlerFunc) http.HandlerFunc {
	if path == "/" {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				NotFoundHandler(w, r)
				return
			}

			handler(w, r)
		}
	}

	return handler
}
