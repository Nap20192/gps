package middleware

import "net/http"

type Middleware func(next http.HandlerFunc) http.HandlerFunc

func CreateMiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
