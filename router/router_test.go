package router

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEngine(t *testing.T) {
	router := New()
	assert.Equal(t, "/", router.basePath)
	assert.Equal(t, router.engine, router)
	assert.Empty(t, router.Handlers)
}

func TestAddRoute(t *testing.T) {
	router := New()
	router.addRoute("/add", HandlersChain{func(_ *Context) {}})
	router.addRoute("/delete", HandlersChain{func(_ *Context) {}})
	router.addRoute("/hello", HandlersChain{func(_ *Context) {}})
	assert.Len(t, router.tree.children, 3)
}

func TestAddRouteFails(t *testing.T) {
	router := New()
	assert.Panics(t, func() { router.addRoute("a", HandlersChain{func(_ *Context) {}}) })
	assert.Panics(t, func() { router.addRoute("/", HandlersChain{}) })

	router.addRoute("/post", HandlersChain{func(_ *Context) {}})
	assert.Panics(t, func() {
		router.addRoute("/post", HandlersChain{func(_ *Context) {}})
	})
}

// func TestCreateDefaultRouter(t *testing.T) {
// 	router := Default()
// 	assert.Len(t, router.Handlers, 2)
// }

func TestNoRouteWithoutGlobalHandlers(t *testing.T) {
	var middleware0 HandlerFunc = func(c *Context) {}
	var middleware1 HandlerFunc = func(c *Context) {}

	router := New()

	router.NoRoute(middleware0)
	assert.Nil(t, router.Handlers)
	assert.Len(t, router.noRoute, 1)
	assert.Len(t, router.allNoRoute, 1)
	compareFunc(t, router.noRoute[0], middleware0)
	compareFunc(t, router.allNoRoute[0], middleware0)

	router.NoRoute(middleware1, middleware0)
	assert.Len(t, router.noRoute, 2)
	assert.Len(t, router.allNoRoute, 2)
	compareFunc(t, router.noRoute[0], middleware1)
	compareFunc(t, router.allNoRoute[0], middleware1)
	compareFunc(t, router.noRoute[1], middleware0)
	compareFunc(t, router.allNoRoute[1], middleware0)
}

func TestNoRouteWithGlobalHandlers(t *testing.T) {
	var middleware0 HandlerFunc = func(c *Context) {}
	var middleware1 HandlerFunc = func(c *Context) {}
	var middleware2 HandlerFunc = func(c *Context) {}

	router := New()
	router.Use(middleware2)

	router.NoRoute(middleware0)
	assert.Len(t, router.allNoRoute, 2)
	assert.Len(t, router.Handlers, 1)
	assert.Len(t, router.noRoute, 1)

	compareFunc(t, router.Handlers[0], middleware2)
	compareFunc(t, router.noRoute[0], middleware0)
	compareFunc(t, router.allNoRoute[0], middleware2)
	compareFunc(t, router.allNoRoute[1], middleware0)

	router.Use(middleware1)
	assert.Len(t, router.allNoRoute, 3)
	assert.Len(t, router.Handlers, 2)
	assert.Len(t, router.noRoute, 1)

	compareFunc(t, router.Handlers[0], middleware2)
	compareFunc(t, router.Handlers[1], middleware1)
	compareFunc(t, router.noRoute[0], middleware0)
	compareFunc(t, router.allNoRoute[0], middleware2)
	compareFunc(t, router.allNoRoute[1], middleware1)
	compareFunc(t, router.allNoRoute[2], middleware0)
}

func compareFunc(t *testing.T, a, b interface{}) {
	sf1 := reflect.ValueOf(a)
	sf2 := reflect.ValueOf(b)
	if sf1.Pointer() != sf2.Pointer() {
		t.Error("different functions")
	}
}

func TestListOfRoutes(t *testing.T) {
	router := New()
	router.Handle("/favicon.ico", handlerTest1)
	router.Handle("/", handlerTest1)
	group := router.Group("/users")
	{
		group.Handle("/", handlerTest2)
		group.Handle("/:id", handlerTest1)
	}
	list := router.Routes()

	assert.Len(t, list, 4)
	assertRoutePresent(t, list, RouteInfo{
		Path:    "/favicon.ico",
		Handler: "^(.*/vendor/)?github.com/sleep2death/hexcore/router.handlerTest1$",
	})
	assertRoutePresent(t, list, RouteInfo{
		Path:    "/",
		Handler: "^(.*/vendor/)?github.com/sleep2death/hexcore/router.handlerTest1$",
	})
	assertRoutePresent(t, list, RouteInfo{
		Path:    "/users/",
		Handler: "^(.*/vendor/)?github.com/sleep2death/hexcore/router.handlerTest2$",
	})
	assertRoutePresent(t, list, RouteInfo{
		Path:    "/users/:id",
		Handler: "^(.*/vendor/)?github.com/sleep2death/hexcore/router.handlerTest1$",
	})
}

func assertRoutePresent(t *testing.T, gotRoutes RoutesInfo, wantRoute RouteInfo) {
	for _, gotRoute := range gotRoutes {
		if gotRoute.Path == wantRoute.Path {
			assert.Regexp(t, wantRoute.Handler, gotRoute.Handler)
			return
		}
	}
	t.Errorf("route not found: %v", wantRoute)
}

func handlerTest1(c *Context) {}
func handlerTest2(c *Context) {}
