package swift

type StaticServer struct {
	staticDir string
	path      string
	spa       bool
}

func (r *Swift) RootStaticServer(staticDir string, spa bool) {
	r.context.staticServer = &StaticServer{
		staticDir: staticDir,
		path:      "/",
		spa:       spa,
	}
}
