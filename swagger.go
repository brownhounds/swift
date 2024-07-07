package swift

type SwaggerServer struct {
	staticDir string
	path      string
	serve     bool
}

func (r *Swift) SwaggerStaticServer(staticDir, path string) {
	r.context.swagger = &SwaggerServer{
		staticDir: staticDir,
		path:      path,
		serve:     false,
	}
}

func (r *Swift) SwaggerServe(serve bool) {
	r.context.swagger.serve = serve
}
