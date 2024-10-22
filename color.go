package zen

import (
	"bytes"
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"io"
	"os"

	_ "github.com/mattn/go-isatty"
)

type (
	inner func(interface{}, []string, *Color) string
)

const Rd = "31"

var (
	red = outer(Rd)
)

func outer(n string) inner {
	return func(msg interface{}, styles []string, c *Color) string {
		if c.disabled {
			return fmt.Sprintf("%v", msg)
		}

		b := new(bytes.Buffer)
		b.WriteString("\x1b[")
		b.WriteString(n)
		for _, s := range styles {
			b.WriteString(";")
			b.WriteString(s)
		}
		b.WriteString("m")
		return fmt.Sprintf("%s%v\x1b[0m", b.String(), msg)
	}
}

type Color struct {
	op       io.Writer
	disabled bool
}

func NewColor() (c *Color) {
	c = new(Color)
	c.setOutput(colorable.NewColorableStdout())
	return
}

func (c *Color) output() io.Writer {
	return c.op
}

func (c *Color) setOutput(w io.Writer) {
	c.op = w
	if w, ok := w.(*os.File); !ok || !isatty.IsTerminal(w.Fd()) {
		c.disabled = true
	}
}

func (c *Color) printF(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(c.op, format, args...)
}

func (c *Color) red(msg interface{}, styles ...string) string {
	return red(msg, styles, c)
}
