package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"regexp"
)

var validRequestID = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)

// RequestID reads or generates a unique request identifier and sets it on the response.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" || len(id) > 64 || !validRequestID.MatchString(id) {
			id = generateID()
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
