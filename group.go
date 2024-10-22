package zen

import (
	"net/http"
)

type Group struct {
	prefix      string
	middlewares []MiddlewareFunc
	zen         *Zen
}

// GET registers a new GET route for a path with matching handler in the router
func (g *Group) GET(pattern string, handlerFunc HandlerFunc) {
	g.Add(http.MethodGet, pattern, handlerFunc)
}

// POST registers a new POST route for a path with matching handler in the router
func (g *Group) POST(pattern string, handlerFunc HandlerFunc) {
	g.Add(http.MethodPost, pattern, handlerFunc)
}

// PUT registers a new POST route for a path with matching handler in the router
func (g *Group) PUT(pattern string, handlerFunc HandlerFunc) {
	g.Add(http.MethodPut, pattern, handlerFunc)
}

// HEAD registers a new HEAD route for a path with matching handler in the router
func (g *Group) HEAD(pattern string, handlerFunc HandlerFunc) {
	g.Add(http.MethodHead, pattern, handlerFunc)
}

// OPTIONS registers a new OPTIONS route for a path with matching handler in the router
func (g *Group) OPTIONS(pattern string, handlerFunc HandlerFunc) {
	g.Add(http.MethodOptions, pattern, handlerFunc)
}

// PATCH registers a new PATCH route for a path with matching handler in the router
func (g *Group) PATCH(pattern string, handlerFunc HandlerFunc) {
	g.Add(http.MethodPatch, pattern, handlerFunc)
}

// DELETE registers a new DELETE route for a path with matching handler in the router
func (g *Group) DELETE(pattern string, handlerFunc HandlerFunc) {
	g.Add(http.MethodDelete, pattern, handlerFunc)
}

// TRACE registers a new TRACE route for a path with matching handler in the router
func (g *Group) TRACE(pattern string, handlerFunc HandlerFunc) {
	g.Add(http.MethodTrace, pattern, handlerFunc)
}

func (g *Group) Add(method string, pattern string, handlerFunc HandlerFunc) {
	g.addRoute(method, pattern, handlerFunc)
}

func (g *Group) addRoute(method string, pattern string, handler HandlerFunc) {
	p := g.prefix + pattern
	g.zen.router.Add(method, p, handler)
}
