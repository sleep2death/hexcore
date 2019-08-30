package router

import (
	"github.com/sleep2death/hexcore/actions"
)

// ActionType - which is used for group actions
type ActionType string

const (
	// Card action type
	Card ActionType = "Card"
	// Battle action type
	Battle ActionType = "Battle"
	// Normal action type
	Normal ActionType = "Normal"
)

// Handle is a function that can be registered to a route to handle Action
type Handle func(Params) actions.Action

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

// Router can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	trees map[string]*node

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

	// Configurable Handler which is called when no matching route is found.
	NotFound actions.Action
}

// New returns a new initialized Router.
// Path auto-correction, including trailing slashes, is enabled by default.
func New() *Router {
	return &Router{
		RedirectTrailingSlash: true,
		RedirectFixedPath:     true,
	}
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) Handle(actionType ActionType, path string, handle Handle) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	group := string(actionType)

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	root := r.trees[group]
	if root == nil {
		root = new(node)
		r.trees[string(actionType)] = root
	}

	root.addRoute(path, handle)
}

// Serve - get the Action by given path
func (r *Router) Serve(actionType ActionType, path string) actions.Action {
	group := string(actionType)
	if root := r.trees[group]; root != nil {
		if handle, ps, tsr := root.getValue(path); handle != nil {
			return handle(ps)
		} else if path != "/" {
			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					path = path[:len(path)-1]
				} else {
					path = path + "/"
				}
				return r.Serve(actionType, path)
			}

			// Try to fix the request path
			if r.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					r.RedirectTrailingSlash,
				)
				if found {
					path = string(fixedPath)
					return r.Serve(actionType, path)
				}
			}
		}
	}

	if r.NotFound != nil {
		return r.NotFound
	}

	return nil
}
