package middlewares

import (
	"fmt"
	"github.com/jeffcail/zen"
)

func Demo() zen.HandlerFunc {
	return func(c *zen.Context) {
		fmt.Println("demo")

		c.Set("example", "123456")

		c.Next()

	}
}
