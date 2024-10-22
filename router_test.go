package zen

import (
	"fmt"
	"reflect"
	"testing"
)

func newTesRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	r := newTesRouter()
	ok := reflect.DeepEqual(r.parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(r.parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(r.parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTesRouter()
	n, ps := r.find("GET", "/hello/zen")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "zen" {
		t.Fatal("name should be equal to 'zen'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])

}
