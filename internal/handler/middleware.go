package handler

import (
	"fmt"
	"net/http"
)

func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("panic middleware", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered", err)
				newErrorResponse(w, `Internal server error`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
