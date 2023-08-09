package serverutils

import (
	"net/http"
	"time"
)

var defender = NewDefender(5, 1*time.Second, 5*time.Minute)

func Antibrut(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if client, ok := defender.Client(r.RemoteAddr); ok && client.Banned() {
			http.Error(w, "Forbidden", 403)
			return
		}
		defender.Inc(r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
