package os

import (
	"os"
)

func init() {
	Exit = os.Exit
}

// Exit will use os.Exit and can be overridden for testing.
var Exit func(code int)
