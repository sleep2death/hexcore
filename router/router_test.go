package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	r := New()
	r.Handle("GET", "/Hello/:name", func(ps Params) {
		assert.Equal(t, "World", ps.ByName("name"))
		// t.Logf("Hello, %s.", ps.ByName("name"))
	})

	path := "/Hello/World"
	if root := r.trees["GET"]; root != nil {
		handle, ps, _ := root.getValue(path)
		handle(ps)
	}
}
