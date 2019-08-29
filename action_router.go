package hexcore

// NewRouter returns a new router instance.
func NewRouter() *Router {
	return &Router{namedRoutes: make(map[string]*Route)}
}

// Router of all the actions
type Router struct {
	ActionNotFound   Action
	ActionNotAllowed Action

	namedRoutes map[string]*Route
}

// Route of the action
type Route struct {
}
