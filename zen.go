package zen

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
)

type HandlerFunc func(c *Context)

type MiddlewareFunc func(c *Context) HandlerFunc

type Zen struct {
	zenMutex sync.RWMutex
	color    *Color
	Server   *http.Server
	Listener net.Listener
	router   *router
	pool     sync.Pool

	HideBanner bool

	middlewares []HandlerFunc

	groups []*Group

	ListenerNetWork string
}

// New create an instance of Zen
func New() (z *Zen) {
	z = &Zen{
		color:           NewColor(),
		Server:          new(http.Server),
		ListenerNetWork: "tcp",
		router:          newRouter(),
	}
	z.Server.Handler = z
	return
}

func (z *Zen) newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		zen:    z,
		r:      r,
		w:      w,
		path:   r.URL.Path,
		method: r.Method,
		mx:     -1,
	}
}

// GET registers a new GET route for a path with matching handler in the router
func (z *Zen) GET(pattern string, handlerFunc HandlerFunc) {
	z.router.Add(http.MethodGet, pattern, handlerFunc)
}

// POST registers a new POST route for a path with matching handler in the router
func (z *Zen) POST(pattern string, handlerFunc HandlerFunc) {
	z.router.Add(http.MethodPost, pattern, handlerFunc)
}

// PUT registers a new PUT route for a path with matching handler in the router
func (z *Zen) PUT(pattern string, handlerFunc HandlerFunc) {
	z.router.Add(http.MethodPut, pattern, handlerFunc)
}

// HEAD registers a new HEAD route for a path with matching handler in the router
func (z *Zen) HEAD(pattern string, handlerFunc HandlerFunc) {
	z.router.addRoute(http.MethodHead, pattern, handlerFunc)
}

// OPTIONS registers a new OPTIONS route for a path with matching handler in the router
func (z *Zen) OPTIONS(pattern string, handlerFunc HandlerFunc) {
	z.router.addRoute(http.MethodOptions, pattern, handlerFunc)
}

// PATCH registers a new PATCH route for a path with matching handler in the router
func (z *Zen) PATCH(pattern string, handlerFunc HandlerFunc) {
	z.router.addRoute(http.MethodPatch, pattern, handlerFunc)
}

// DELETE registers a new DELETE route for a path with matching handler in the router
func (z *Zen) DELETE(pattern string, handlerFunc HandlerFunc) {
	z.router.addRoute(http.MethodDelete, pattern, handlerFunc)
}

// TRACE registers a new TRACE route for a path with matching handler in the router
func (z *Zen) TRACE(pattern string, handlerFunc HandlerFunc) {
	z.router.addRoute(http.MethodTrace, pattern, handlerFunc)
}

// Use is defined to add middleware to the group
func (z *Zen) Use(middlewares ...HandlerFunc) {
	z.middlewares = append(z.middlewares, middlewares...)
}

// Start an http server
func (z *Zen) Start(addr string) (err error) {
	z.zenMutex.Lock()
	defer z.zenMutex.Unlock()
	z.Server.Addr = addr
	if err = z.configureServer(z.Server); err != nil {
		return err
	}
	return z.Server.Serve(z.Listener)
}

// SetHideBanner set banner is hide.
func (z *Zen) SetHideBanner(b bool) {
	z.HideBanner = b
}

// GetHideBanner get hide banner value.
func (z *Zen) GetHideBanner() bool {
	return z.HideBanner
}

func (z *Zen) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//c := z.newContext(w, r)
	//z.router.handle(c)
	var middlewares []HandlerFunc
	for _, group := range z.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := z.newContext(w, r)
	c.middlewares = middlewares
	z.router.handle(c)
}

func (z *Zen) configureServer(s *http.Server) error {

	if !z.HideBanner {
		z.color.printF(banner, z.color.red("v"+version), z.color.red(github))
		z.color.printF(fmt.Sprintf("=> port %s %s\n", z.color.red("[::]"), z.Server.Addr))
	}

	if z.Listener == nil {
		l, err := newListener(s.Addr, z.ListenerNetWork)
		if err != nil {
			return err
		}
		z.Listener = l
	}

	return nil
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func newListener(address, network string) (*tcpKeepAliveListener, error) {
	if network != TcpNetwork && network != TcpNet4work && network != TcpNet6work {
		return nil, ErrInvalidListenerNetwork
	}
	l, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	return &tcpKeepAliveListener{l.(*net.TCPListener)}, nil
}

func (z *Zen) Group(prefix string) *Group {
	ng := &Group{
		prefix: prefix,
		zen:    z,
	}
	z.groups = append(z.groups, ng)
	return ng
}
