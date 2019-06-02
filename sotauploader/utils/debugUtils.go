package utils

import (
	"os"
)

func SetDebugInput(debug bool) {
	if debug {
		os.Setenv("REQUEST_METHOD", "GET")
		os.Setenv("SERVER_PROTOCOL", "HTTP/1.1")
		os.Setenv("HTTP_COOKIE", "CMSSESSIDac42f1aa=v86fvkfsad45qcfe24u2h4c072")
	}
}
