package serverutils

import (
	"net/http"
	"os"
)

func AuthWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		if os.Getenv("DEV") == "true" {
			next.ServeHTTP(w, r)
			return
		}
		username, password, authOK := r.BasicAuth()
		if authOK == false {
			http.Error(w, "Not authorized", 401)
			return
		}

		if username != "admin" || password != "9qfpL6Mrwyax7MZvuYXz" {
			http.Error(w, "Not authorized", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
