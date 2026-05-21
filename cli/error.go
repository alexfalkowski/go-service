package cli

import "strconv"

type shutdownError int

func (e shutdownError) Error() string {
	return "shutdown failed with exit code " + strconv.Itoa(e.code())
}

func (e shutdownError) code() int {
	return int(e)
}
