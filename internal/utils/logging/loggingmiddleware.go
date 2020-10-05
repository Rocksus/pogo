package logging

import (
	"net/http"
	"time"

	"github.com/nickylogan/go-log"
)

// Middleware handles the calculation of processing time the bot takes
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmStart := time.Now()
		next.ServeHTTP(w, r)
		tmEnd := time.Now()
		log.Infoln("Process Time:", tmEnd.Sub(tmStart))
	})
}
