package swift

import "net/http"

type Group struct {
	swift       *Swift
	path        string
	middlewares []Middleware
}

func (r *Swift) Group(path string) *Group {
	p := BuildAndValidatePath(path)
	return &Group{
		swift:       r,
		path:        p,
		middlewares: make([]Middleware, 0),
	}
}

func (g *Group) Group(path string) *Group {
	p := BuildAndValidatePath(path)
	return &Group{
		swift:       g.swift,
		path:        g.path + p,
		middlewares: g.middlewares,
	}
}

func (g *Group) Middleware(m ...Middleware) *Group {
	g.middlewares = append(g.middlewares, m...)
	return g
}

func (g *Group) Get(path string, handler http.HandlerFunc) *HandlerValue {
	p := BuildAndValidatePath(path)
	return g.swift.MakeHandler(http.MethodGet, g.path+p, handler, g)
}

func (g *Group) Post(path string, handler http.HandlerFunc) *HandlerValue {
	p := BuildAndValidatePath(path)
	return g.swift.MakeHandler(http.MethodPost, g.path+p, handler, nil)
}

func (g *Group) Put(path string, handler http.HandlerFunc) *HandlerValue {
	p := BuildAndValidatePath(path)
	return g.swift.MakeHandler(http.MethodPut, g.path+p, handler, nil)
}

func (g *Group) Patch(path string, handler http.HandlerFunc) *HandlerValue {
	p := BuildAndValidatePath(path)
	return g.swift.MakeHandler(http.MethodPatch, g.path+p, handler, nil)
}

func (g *Group) Delete(path string, handler http.HandlerFunc) *HandlerValue {
	p := BuildAndValidatePath(path)
	return g.swift.MakeHandler(http.MethodDelete, g.path+p, handler, nil)
}
