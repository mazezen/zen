package middlewares

import (
	"github.com/jeffcail/zen"
	"log"
	"time"
)

func Demo2() zen.HandlerFunc {
	return func(c *zen.Context) {
		t := time.Now()

		c.Set("example2", "2222222")

		c.Next()

		latency := time.Since(t)
		log.Println(latency)

	}
}
