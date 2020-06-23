package server

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

var mappings = make(map[string]string)

// GetMappings gets all mappings
func GetMappings() map[string]string {
	return mappings
}

// CreateRoute creates a new route
func CreateRoute(hostname string, backend string) {
	hostname = strings.ToLower(hostname)
	logrus.WithFields(logrus.Fields{
		"hostname": hostname,
		"upstream": backend,
	}).Info("route created")
	mappings[hostname] = backend
}

// ReadRoute reads a route by an address
func ReadRoute(address string) (string, error) {
	addressParts := strings.Split(address, "\x00")
	hostname := strings.ToLower(addressParts[0])

	if route, exists := mappings[hostname]; exists {
		return route, nil
	}
	return "", errors.New("no matching route")
}

// DeleteRoute deletes a route by its hostname
func DeleteRoute(hostname string) {
	hostname = strings.ToLower(hostname)
	logrus.WithFields(logrus.Fields{
		"hostname": hostname,
	}).Info("route deleted")
	delete(mappings, hostname)
}
