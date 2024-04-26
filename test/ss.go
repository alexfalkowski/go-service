package test

import (
	"os/exec"
	"strings"

	"github.com/alexfalkowski/go-service/runtime"
)

// KillPort using ss.
func KillPort(port string) {
	c := "sudo ss --kill state listening src :" + port
	s := strings.Split(c, " ")
	cmd := exec.Command(s[0], s[1:]...) //nolint:gosec

	runtime.Must(cmd.Start())

	runtime.Must(cmd.Wait())
}
