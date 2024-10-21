package main

import (
	"github.com/jeffcail/zen"
	"net/http"
)

// curl http://127.0.0.1:8888/
// curl http://127.0.0.1:8888/ping
func main() {
	z := zen.New()

	z.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
		w.WriteHeader(http.StatusOK)
	})

	z.Start(":8888")
}
