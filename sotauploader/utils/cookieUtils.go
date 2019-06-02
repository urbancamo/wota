package utils

import (
	"net/http"
	"strings"
)

func FindAuthCookie(cookies []*http.Cookie) string {
	for _, cookie := range cookies {
		if strings.Index(cookie.Name, "CMSSESSID") >= 0 {
			return cookie.Value
		}
	}
	return ""
}
