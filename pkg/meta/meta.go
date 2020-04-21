package meta

import (
	"fmt"
)

var branch string
var commit string
var version string

func Name() string {
	return fmt.Sprintf("tq %s", Version())
}

func Version() string {
	if len(branch) > 0 && len(commit) > 0 {
		return fmt.Sprintf("%s@%s", commit, branch)
	} else if len(version) > 0 {
		return fmt.Sprintf("%s", version)
	} else {
		return "unknown"
	}
}
