package router

import (
	"errors"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

type TestWriter struct {
}

func (w *TestWriter) Write(p []byte) (int, error) {
	return 0, nil
}

func (w *TestWriter) Close() error {
	return nil
}

func performRequest(r *Engine, path string) {
	msg := &any.Any{TypeUrl: path}
	buf, _ := proto.Marshal(msg)
	r.Serve(buf, &TestWriter{})
}

func TestMiddlewareGeneralCase(t *testing.T) {
	signature := ""
	router := New()
	router.Use(func(c *Context) {
		signature += "A"
		c.Next()
		signature += "B"
	})
	router.Use(func(c *Context) {
		signature += "C"
	})
	router.Handle("/", func(c *Context) {
		signature += "D"
	})
	router.NoRoute(func(c *Context) {
		signature += " X "
	})
	// RUN
	performRequest(router, "/")

	// TEST
	assert.Equal(t, "ACDB", signature)
}

func TestMiddlewareNoRoute(t *testing.T) {
	signature := ""
	router := New()
	router.Use(func(c *Context) {
		signature += "A"
		c.Next()
		signature += "B"
	})
	router.Use(func(c *Context) {
		signature += "C"
		c.Next()
		c.Next()
		c.Next()
		c.Next()
		signature += "D"
	})
	router.NoRoute(func(c *Context) {
		signature += "E"
		c.Next()
		signature += "F"
	}, func(c *Context) {
		signature += "G"
		c.Next()
		signature += "H"
	})
	// RUN
	performRequest(router, "/")

	// TEST
	assert.Equal(t, "ACEGHFDB", signature)
}
func TestMiddlewareAbort(t *testing.T) {
	signature := ""
	router := New()
	router.Use(func(c *Context) {
		signature += "A"
	})
	router.Use(func(c *Context) {
		signature += "C"
		c.Abort()
		c.Next()
		signature += "D"
	})
	router.Handle("/", func(c *Context) {
		signature += " X "
		c.Next()
		signature += " XX "
	})

	// RUN
	performRequest(router, "/")

	// TEST
	assert.Equal(t, "ACD", signature)
}

func TestMiddlewareAbortHandlersChainAndNext(t *testing.T) {
	signature := ""
	router := New()
	router.Use(func(c *Context) {
		signature += "A"
		c.Next()
		c.Abort()
		signature += "B"

	})
	router.Handle("/", func(c *Context) {
		signature += "C"
		c.Next()
	})
	// RUN
	performRequest(router, "/")

	// TEST
	assert.Equal(t, "ACB", signature)
}

// TestFailHandlersChain - ensure that Fail interrupt used middleware in fifo order as
// as well as Abort
func TestMiddlewareFailHandlersChain(t *testing.T) {
	// SETUP
	signature := ""
	router := New()
	router.Use(func(context *Context) {
		signature += "A"
		context.AbortWithError(errors.New("foo")) // nolint: errcheck
	})
	router.Use(func(context *Context) {
		signature += "B"
		context.Next()
		signature += "C"
	})
	// RUN
	performRequest(router, "/")

	// TEST
	assert.Equal(t, "A", signature)
}
