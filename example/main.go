package main

import (
	"fmt"
	"github.com/mazezen/zen"
	"github.com/mazezen/zen/example/middlewares"
	"log"
	"net/http"
	"strconv"
)

//// curl http://127.0.0.1:8888/
//// curl http://127.0.0.1:8888/ping
//func main() {
//	z := zen.New()
//
//	z.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte("pong"))
//		w.WriteHeader(http.StatusOK)
//	})
//
//	z.Start(":8888")
//}

func main() {
	z := zen.New()

	z.Use(middlewares.Demo())

	// curl http://127.0.0.1:8888
	z.GET("/", func(c *zen.Context) {
		c.HTML(http.StatusOK, "<h1>hello world</h1>")
	})

	// curl http://127.0.0.1:8888/query?name=zen&email=zen@zen.com
	z.GET("/query", func(c *zen.Context) {
		c.JSON(http.StatusOK, zen.H{
			"name":  c.Query("name"),
			"email": c.Query("email"),
		})
	})

	// curl http://127.0.0.1:8888/hello/zen
	z.GET("/hello/:name", func(c *zen.Context) {
		fmt.Println(c.Get("example"))
		c.String(http.StatusOK, "hello %s, you are beautiful!", c.Param("name"))
	}, middlewares.Demo())

	// curl http://127.0.0.1:8888/assets/css/demo.css
	z.GET("/assets/*filepath", func(c *zen.Context) {
		c.String(http.StatusOK, c.Param("filepath"))
	})

	// curl http://127.0.0.1:8888/put/12
	z.PUT("/put/:id", func(c *zen.Context) {
		var id = 10
		v := c.Param("id")
		av, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("convert %s to int error: %v", v, err)
		}
		id = av
		c.JSON(http.StatusOK, zen.H{
			"message": "update id value is success",
			"id":      id,
		})
	})

	// curl http://127.0.0.1:8888/delete/12
	z.DELETE("/delete/:id", func(c *zen.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, zen.H{
			"message": "delete id value is " + id,
			"status":  "success",
		})
	})

	// curl "http://127.0.0.1:8888/login" -X POST -d 'name=zen&email=zen@gmail.com'
	z.POST("/login", func(c *zen.Context) {
		type Data struct {
			Name    string            `json:"name"`
			Email   string            `json:"email"`
			DataMap map[string]string `json:"data_map"`
		}
		d := &Data{
			Name:  c.FormValue("name"),
			Email: c.FormValue("email"),
			DataMap: map[string]string{
				"name":  c.FormValue("name"),
				"email": c.FormValue("email"),
			},
		}

		c.JSON(http.StatusOK, d)
	})

	// curl http://127.0.0.1:8888/v2/demo
	v2 := z.Group("/v2")
	v2.Use(middlewares.Demo())
	v2.GET("/demo", func(c *zen.Context) {
		c.String(http.StatusOK, "demo router group!")
		fmt.Println(c.Get("example"))
		fmt.Println(c.Get("example2"))

	}, middlewares.Demo(), middlewares.Demo2())

	// curl "http://127.0.0.1:8888/v2/post" -X POST -d 'name=zen&email=zen@gmail.com'
	v2.POST("/post", func(c *zen.Context) {
		type Data struct {
			Name    string            `json:"name"`
			Email   string            `json:"email"`
			DataMap map[string]string `json:"data_map"`
		}
		d := &Data{
			Name:  c.FormValue("name"),
			Email: c.FormValue("email"),
			DataMap: map[string]string{
				"name":  c.FormValue("name"),
				"email": c.FormValue("email"),
			},
		}

		value := c.Get("example")
		log.Println(value)

		c.JSON(http.StatusOK, d)
	})

	if err := z.Start(":8888"); err != nil {
		log.Printf("zen start error: %v", err)
	}

}
