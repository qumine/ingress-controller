package build

import (
	"fmt"
)

// nolint: gochecknoglobals
var (
	version = "dev"
	commit  = ""
	date    = ""
)

func GetVersion() string {
	result := version
	if commit != "" {
		result = fmt.Sprintf("%s, commit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s, built at: %s", result, date)
	}
	return result
}
