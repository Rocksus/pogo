package logging

import (
	"log"
	"net/http"
	"time"
)

// Middleware handles the calculation of processing time the bot takes
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmStart := time.Now()
		next.ServeHTTP(w, r)
		tmEnd := time.Now()
		log.Print("[Chatbot Call] Process Time: ", tmEnd.Sub(tmStart))
	})
}
