package serverutils

import "net/http"

type MiddlewareFunc = func(handler http.Handler) http.Handler
type HandlerFunc = func(w http.ResponseWriter, r *http.Request)

func Compose(handler http.Handler, middlewares ...MiddlewareFunc) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
