package main

import (
	"fmt"
	"net/http"
)

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("testing from middleware")
		authHeader := r.Header.Get("Authorization-Token")
		if authHeader != "1234" {
			http.Error(w, "Unauthorized: incorrect auth Header", http.StatusUnauthorized)
		}
		next.ServeHTTP(w, r)
	})
}
