package zen

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Zen struct {
	zenMutex sync.RWMutex
	color    *Color
	Server   *http.Server
	Listener net.Listener
	router   map[string]HandlerFunc

	HideBanner bool

	ListenerNetWork string
}

// New create an instance of Zen
func New() (z *Zen) {
	z = &Zen{
		color:           NewColor(),
		Server:          new(http.Server),
		ListenerNetWork: "tcp",
	}
	z.Server.Handler = z
	z.router = map[string]HandlerFunc{}
	return
}

func (z *Zen) addRoute(method, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	z.router[key] = handler
}

// GET registers a new GET route for a path with matching handler in the router
func (z *Zen) GET(pattern string, handler HandlerFunc) {
	z.addRoute("GET", pattern, handler)
}

// POST registers a new POST route for a path with matching handler in the router
func (z *Zen) POST(pattern string, handler HandlerFunc) {
	z.addRoute("POST", pattern, handler)
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
	key := r.Method + "-" + r.URL.Path
	if handler, ok := z.router[key]; ok {
		handler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not Found: %s\n", r.URL)
	}
}

func (z *Zen) configureServer(s *http.Server) error {

	if !z.HideBanner {
		z.color.printF(banner, z.color.red("v"+version), z.color.red(github))
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

var ErrInvalidListenerNetwork = errors.New("invalid listener network")

const (
	TcpNetwork  = "tcp"
	TcpNet4work = "tcp4"
	TcpNet6work = "tcp6"
)

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
