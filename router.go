package zen

import (
	"net/http"
	"strings"
)

type router struct {
	tree     map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		tree:     make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
func (r *router) parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) Add(method, path string, handler HandlerFunc) {
	r.addRoute(method, path, handler)
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := r.parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.tree[method]
	if !ok {
		r.tree[method] = &node{}
	}
	r.tree[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) Find(method, path string) (*node, map[string]string) {
	return r.getRoute(method, path)
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := r.parsePattern(path)
	params := make(map[string]string)
	root, ok := r.tree[method]

	if !ok {
		return nil, nil
	}

	n := root.find(searchParts, 0)

	if n != nil {
		parts := r.parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.Find(c.method, c.path)

	if n != nil {
		key := c.method + "-" + n.pattern
		c.params = params
		c.middlewares = append(c.middlewares, r.handlers[key])
	} else {
		c.middlewares = append(c.middlewares, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.path)
		})
	}
	c.Next()
}

type node struct {
	pattern     string
	part        string
	children    []*node
	isPrecision bool
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isPrecision {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isPrecision {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isPrecision: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) find(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.find(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
