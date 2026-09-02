package telegram

type Handler func(Context) error

type Route struct {
	Match   func(Context) bool
	Handler Handler
}

type Router struct {
	routes []Route
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Register(routes ...Route) {
	r.routes = append(r.routes, routes...)
}

func (r *Router) Handle(ctx Context) error {
	for _, route := range r.routes {
		if route.Match(ctx) {
			return route.Handler(ctx)
		}
	}

	return nil
}
