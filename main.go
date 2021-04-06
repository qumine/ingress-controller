package main

import (
	"os"

	"github.com/qumine/ingress-controller/cmd"
	"github.com/sirupsen/logrus"
)

// nolint: gochecknoglobals
var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

//go:generate make generate-bindata
func init() {
	// Set
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	cmd.Execute()
}
