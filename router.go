package hexcore

import "sync"

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
	Method      string
	Path        string
	Handler     string
	HandlerFunc HandlerFunc
}

// RoutesInfo defines a RouteInfo array.
type RoutesInfo []RouteInfo

// Engine -
type Engine struct {
	RouterGroup

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	allNoRoute  HandlersChain
	allNoMethod HandlersChain
	noRoute     HandlersChain
	noMethod    HandlersChain

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
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		RedirectTrailingSlash: true,
		RedirectFixedPath:     true,
		tree:                  new(node),
	}
	engine.tree.fullPath = "/"

	engine.RouterGroup.engine = engine
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	return engine
}

// Default returns an Engine instance with the Logger and Recovery middleware already attached.
func Default() *Engine {
	// debugPrintWARNINGDefault()
	engine := New()
	// engine.Use(Logger(), Recovery())
	// TODO: logger and recovery
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
			Method:      "",
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

// Serve conforms to the http.Handler interface.
func (engine *Engine) Serve(path string) {
	c := engine.pool.Get().(*Context)
	// c.writermem.reset(w)
	c.Path = path
	c.reset()

	engine.handleRequest(c)

	engine.pool.Put(c)
}

func (engine *Engine) handleRequest(c *Context) {
	rPath := c.Path
	rPath = cleanPath(rPath)
	unescape := false

	// Find root of the tree for the given HTTP method
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
	if rPath != "/" {
		if value.tsr && engine.RedirectTrailingSlash {
			redirectTrailingSlash(c)
			return
		}
		if engine.RedirectFixedPath && redirectFixedPath(c, root, engine.RedirectFixedPath) {
			return
		}
	}
}
