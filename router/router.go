package router

import (
	"sync"
)

// HandlerFunc defines the handler used by gin middleware as return value.
type HandlerFunc func(*Context)

// HandlersChain defines a HandlerFunc array.
type HandlersChain []HandlerFunc

// Last returns the last handler in the chain. ie. the last handler is the main own.
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

// RouteInfo represents a request route's specification which contains method and path and its handler.
type RouteInfo struct {
	Path        string
	Handler     string
	HandlerFunc HandlerFunc
}

// RoutesInfo defines a RouteInfo array.
type RoutesInfo []RouteInfo

// Engine -
type Engine struct {
	RGroup

	allNoRoute HandlersChain
	noRoute    HandlersChain

	pool sync.Pool
	tree *node
}

var _ IRouter = &Engine{}

// New returns a new blank Engine instance without any middleware attached.
// By default the configuration is:
// - RedirectTrailingSlash:  true
// - RedirectFixedPath:      false
// - HandleMethodNotAllowed: false
// - ForwardedByClientIP:    true
// - UseRawPath:             false
// - UnescapePathValues:     true
func New() *Engine {
	/// debugPrintWARNINGNew()
	engine := &Engine{
		RGroup: RGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		tree: new(node),
	}
	engine.tree.fullPath = "/"

	engine.RGroup.engine = engine
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	return engine
}

// Default returns an Engine instance with the Logger and Recovery middleware already attached.
func Default() *Engine {
	// debugPrintWARNINGDefault()
	engine := New()
	engine.Use(Logger(), Recovery())
	// TODO: logger
	return engine
}

func (engine *Engine) allocateContext() *Context {
	return &Context{engine: engine}
}

func (engine *Engine) addRoute(path string, handlers HandlersChain) {
	assert1(path[0] == '/', "path must begin with '/'")
	// assert1(method != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")

	// debugPrintRoute(method, path, handlers)
	root := engine.tree
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		// engine.trees = append(engine.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, handlers)
}

// Routes returns a slice of registered routes, including some useful information, such as:
// the http method, path and the handler name.
func (engine *Engine) Routes() (routes RoutesInfo) {
	routes = iterate("", routes, engine.tree)
	return routes
}

func iterate(path string, routes RoutesInfo, root *node) RoutesInfo {
	path += root.path
	if len(root.handlers) > 0 {
		handlerFunc := root.handlers.Last()
		routes = append(routes, RouteInfo{
			Path:        path,
			Handler:     nameOfFunction(handlerFunc),
			HandlerFunc: handlerFunc,
		})
	}
	for _, child := range root.children {
		routes = iterate(path, routes, child)
	}
	return routes
}

// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	engine.noRoute = handlers
	engine.rebuild404Handlers()
}

// Use attaches a global middleware to the router. ie. the middleware attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {
	engine.RGroup.Use(middleware...)
	engine.rebuild404Handlers()
	// engine.rebuild405Handlers()
	return engine
}

func (engine *Engine) rebuild404Handlers() {
	engine.allNoRoute = engine.combineHandlers(engine.noRoute)
}

// Serve with the given path
func (engine *Engine) Serve(path string) {
	c := engine.pool.Get().(*Context)
	// c.writermem.reset(w)
	c.reset()

	c.Path = path

	engine.handleRequest(c)

	engine.pool.Put(c)
}

func (engine *Engine) handleRequest(c *Context) {
	rPath := c.Path
	rPath = cleanPath(rPath)
	unescape := false

	root := engine.tree
	// Find route in tree
	value := root.getValue(rPath, c.Params, unescape)
	if value.handlers != nil {
		c.handlers = value.handlers
		c.Params = value.params
		c.fullPath = value.fullPath
		c.Next()
		// c.writermem.WriteHeaderNow()
		return
	}

	c.handlers = engine.allNoRoute
	serveError(c)
}

func serveError(c *Context) {
	// TODO: serve error
	c.Next()
}