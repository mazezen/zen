package middlewares

import (
	"github.com/jeffcail/zen"
	"log"
	"time"
)

func Demo() zen.HandlerFunc {
	return func(c *zen.Context) {
		t := time.Now()

		c.Set("example", "123456")

		c.Next()

		latency := time.Since(t)
		log.Println(latency)

	}
}
