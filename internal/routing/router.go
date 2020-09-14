package routing

import (
	"errors"
	"strings"

	"github.com/quhive/qumine-ingress/internal/metrics"
	"github.com/sirupsen/logrus"
)

var routes = make(map[string]Route)

// Add a new route to the router.
func Add(uid string, route Route) {
	if _, ok := routes[uid]; !ok {
		routes[uid] = route
		logrus.WithField("uid", uid).WithField("frontend", route.Frontend).WithField("backend", route.Backend).Info("route created")
		metrics.Routes.Inc()
	} else {
		logrus.WithField("uid", uid).Warn("route already created")
	}
}

// Update an existing route from the router.
func Update(uid string, route Route) {
	if _, ok := routes[uid]; ok {
		routes[uid] = route
		logrus.WithField("uid", uid).WithField("frontend", route.Frontend).WithField("backend", route.Backend).Info("route updated")
	}
}

// Remove an existing route from the router.
func Remove(uid string) {
	if _, ok := routes[uid]; ok {
		delete(routes, uid)
		logrus.WithField("uid", uid).Info("route deleted")
		metrics.Routes.Dec()
	}
}

// FindBackend finds a route by its frontend and returns the backend or throws an error.
func FindBackend(frontend string) (string, error) {
	frontendParts := strings.Split(frontend, "\x00")
	frontend = strings.ToLower(frontendParts[0])

	for _, route := range routes {
		if route.Frontend == frontend {
			return route.Backend, nil
		}
	}
	return "", errors.New("route not found")
}
