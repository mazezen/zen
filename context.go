package zen

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

type H map[string]any

type IContext interface {
	FormValue(name string) string
	Query(name string) string
	Param(name string) string
	SetStatusCode(code int)
	SetHeader(name, value string)
	String(code int, format string, values ...interface{})
	JSON(code int, content interface{})
	HTML(code int, content string)
	Set(key string, value any)
	Get(key string) any
	GetAndDelete(key string) (any, bool)
	Delete(key string)
}

type Context struct {
	r   *http.Request
	w   http.ResponseWriter
	zen *Zen

	// request info
	query  url.Values
	path   string
	method string
	params map[string]string

	// response info
	StatusCode int

	middlewares []HandlerFunc
	mx          int

	store sync.Map
}

// FormValue curl "http://127.0.0.1:8888/login" -X POST -d 'name=zen&email=zen@gmail.com'
func (c *Context) FormValue(name string) string {
	return c.r.FormValue(name)
}

// Query curl http://127.0.0.1:8888/query?name=zen&email=zen@zen.com
func (c *Context) Query(name string) string {
	if c.query == nil {
		c.query = c.r.URL.Query()
	}
	return c.query.Get(name)
}

// Param curl http://127.0.0.1:8888/hello/zen
func (c *Context) Param(name string) string {
	value, _ := c.params[name]
	return value
}

func (c *Context) SetStatusCode(code int) {
	c.StatusCode = code
	c.w.WriteHeader(code)
}

func (c *Context) SetHeader(key, value string) {
	c.w.Header().Set(key, value)
}

// String c.String(http.StatusOK, "hello %s, you are beautiful!", c.Param("name"))
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader(HeaderContentType, HeaderString)
	c.SetStatusCode(code)
	_, _ = c.w.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON Json
//
//	c.Json(http.StatusOK, zen.H{
//				"name":  c.Query("name"),
//				"email": c.Query("email"),
//			})
func (c *Context) JSON(code int, content interface{}) {
	c.SetHeader(HeaderContentType, HeaderJson)
	c.SetStatusCode(code)
	if err := json.NewEncoder(c.w).Encode(content); err != nil {
		http.Error(c.w, err.Error(), 500)
	}
}

// HTML c.HTML(http.StatusOK, "<h1>hello world</h1>")
func (c *Context) HTML(code int, content string) {
	c.SetHeader(HeaderContentType, HeaderHTML)
	c.SetStatusCode(code)
	_, _ = c.w.Write([]byte(content))
}

// Set save data in the context
func (c *Context) Set(key string, value any) {
	c.store.Store(key, value)
}

// Get return the value for given key
func (c *Context) Get(key string) any {
	value, ok := c.store.Load(key)
	if !ok {
		return nil
	}
	return value
}

// GetAndDelete return the value for given key and delete given key
func (c *Context) GetAndDelete(key string) (any, bool) {
	return c.store.LoadAndDelete(key)
}

// Delete delete data from context
func (c *Context) Delete(key string) {
	c.store.Delete(key)
}

func (c *Context) Next() {
	c.mx++
	s := len(c.middlewares)
	for ; c.mx < s; c.mx++ {
		c.middlewares[c.mx](c)
	}
}
