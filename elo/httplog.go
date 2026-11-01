package elo

import (
	"net/http"
	"time"

	"github.com/wlbr/commons/log"
)

func LogHTTPRequests(service string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		next.ServeHTTP(writer, request)
		log.Print("%s request: method=%s path=%s remote=%s duration=%s",
			service, request.Method, request.URL.Path, request.RemoteAddr, time.Since(start))
	})
}
