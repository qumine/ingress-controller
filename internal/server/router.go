package server

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

var routes = make(map[string]string)

// AddRoute adds a new route to the routing map.
func AddRoute(hostname string, backend string) {
	routes[hostname] = backend
	logrus.WithField("hostname", hostname).WithField("upstream", backend).Info("route added")
}

// FindRoute finds a route by its address or throws an error.
func FindRoute(address string) (string, error) {
	addressParts := strings.Split(address, "\x00")
	hostname := strings.ToLower(addressParts[0])

	if route, exists := routes[hostname]; exists {
		return route, nil
	}
	return "", errors.New("route not found")
}

// RemoveRoute removes a route from the routing map.
func RemoveRoute(hostname string) {
	delete(routes, hostname)
	logrus.WithField("hostname", hostname).Info("route removed")
}
