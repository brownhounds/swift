package swift

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func resetContext() {
	ContextValue = Context{
		apiPrefix: "",
	}
}

func headerMiddleware(value string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Group", value)
			next.ServeHTTP(w, r)
		})
	}
}

func TestGroupMiddlewareAppliesToAllMethods(t *testing.T) {
	resetContext()
	server := New()
	group := server.Group("/api").Middleware(headerMiddleware("1"))

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	group.Get("/get", handler)
	group.Post("/post", handler)
	group.Put("/put", handler)
	group.Patch("/patch", handler)
	group.Delete("/delete", handler)

	Boot(server)

	cases := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/get"},
		{http.MethodPost, "/api/post"},
		{http.MethodPut, "/api/put"},
		{http.MethodPatch, "/api/patch"},
		{http.MethodDelete, "/api/delete"},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(tc.method, tc.path, http.NoBody)
		res := httptest.NewRecorder()
		server.serverMux.ServeHTTP(res, req)

		if res.Result().StatusCode != http.StatusOK {
			t.Fatalf("%s %s: expected status %d, got %d", tc.method, tc.path, http.StatusOK, res.Result().StatusCode)
		}

		values := res.Result().Header.Values("X-Group")
		if len(values) != 1 || values[0] != "1" {
			t.Fatalf("%s %s: expected X-Group header [1], got %v", tc.method, tc.path, values)
		}
	}
}

func TestNestedGroupsDoNotLeakMiddlewareToParents(t *testing.T) {
	resetContext()
	server := New()
	base := server.Group("/api").Middleware(headerMiddleware("base"))
	child := base.Group("/v1").Middleware(headerMiddleware("child"))

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	base.Get("/ping", handler)
	child.Get("/ping", handler)

	Boot(server)

	parentReq := httptest.NewRequest(http.MethodGet, "/api/ping", http.NoBody)
	parentRes := httptest.NewRecorder()
	server.serverMux.ServeHTTP(parentRes, parentReq)

	parentValues := parentRes.Result().Header.Values("X-Group")
	if len(parentValues) != 1 || parentValues[0] != "base" {
		t.Fatalf("parent route: expected X-Group header [base], got %v", parentValues)
	}

	childReq := httptest.NewRequest(http.MethodGet, "/api/v1/ping", http.NoBody)
	childRes := httptest.NewRecorder()
	server.serverMux.ServeHTTP(childRes, childReq)

	childValues := childRes.Result().Header.Values("X-Group")
	if len(childValues) != 2 || childValues[0] != "base" || childValues[1] != "child" {
		t.Fatalf("child route: expected X-Group header [base child], got %v", childValues)
	}
}

func TestNestedGroupsAllMethodsHaveMiddleware(t *testing.T) {
	resetContext()
	server := New()
	parent := server.Group("/api").Middleware(headerMiddleware("parent"))
	child := parent.Group("/v1").Middleware(headerMiddleware("child"))

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	parent.Get("/get", handler)
	parent.Post("/post", handler)
	parent.Put("/put", handler)
	parent.Patch("/patch", handler)
	parent.Delete("/delete", handler)

	child.Get("/get", handler)
	child.Post("/post", handler)
	child.Put("/put", handler)
	child.Patch("/patch", handler)
	child.Delete("/delete", handler)

	Boot(server)

	parentCases := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/get"},
		{http.MethodPost, "/api/post"},
		{http.MethodPut, "/api/put"},
		{http.MethodPatch, "/api/patch"},
		{http.MethodDelete, "/api/delete"},
	}

	for _, tc := range parentCases {
		req := httptest.NewRequest(tc.method, tc.path, http.NoBody)
		res := httptest.NewRecorder()
		server.serverMux.ServeHTTP(res, req)

		if res.Result().StatusCode != http.StatusOK {
			t.Fatalf("parent %s %s: expected status %d, got %d", tc.method, tc.path, http.StatusOK, res.Result().StatusCode)
		}

		values := res.Result().Header.Values("X-Group")
		if len(values) != 1 || values[0] != "parent" {
			t.Fatalf("parent %s %s: expected X-Group header [parent], got %v", tc.method, tc.path, values)
		}
	}

	childCases := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/get"},
		{http.MethodPost, "/api/v1/post"},
		{http.MethodPut, "/api/v1/put"},
		{http.MethodPatch, "/api/v1/patch"},
		{http.MethodDelete, "/api/v1/delete"},
	}

	for _, tc := range childCases {
		req := httptest.NewRequest(tc.method, tc.path, http.NoBody)
		res := httptest.NewRecorder()
		server.serverMux.ServeHTTP(res, req)

		if res.Result().StatusCode != http.StatusOK {
			t.Fatalf("child %s %s: expected status %d, got %d", tc.method, tc.path, http.StatusOK, res.Result().StatusCode)
		}

		values := res.Result().Header.Values("X-Group")
		if len(values) != 2 || values[0] != "parent" || values[1] != "child" {
			t.Fatalf("child %s %s: expected X-Group header [parent child], got %v", tc.method, tc.path, values)
		}
	}
}
