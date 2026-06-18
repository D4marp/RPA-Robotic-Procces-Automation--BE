package initenv

import (
	"os"
)

func init() {
	mode := os.Getenv("GIN_MODE")
	if mode != "" && mode != "debug" && mode != "release" && mode != "test" {
		os.Setenv("GIN_MODE", "release")
	}
}
