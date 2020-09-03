package routing

// Route represents the route between a frontend and a backend.
type Route struct {
	Frontend string
	Backend  string
}

// NewRoute creates a new route.
func NewRoute(frontend string, backend string) Route {
	return Route{
		Frontend: frontend,
		Backend:  backend,
	}
}
