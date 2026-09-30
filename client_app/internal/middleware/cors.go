package middleware

import (
	"net/http"
)

func (m *Miiiddeware) Corsmiddwaare(frontorigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resspo http.ResponseWriter, req *http.Request) {

			// //fmt.Println("--######-",frontorigin)

			resspo.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			resspo.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
			resspo.Header().Set("Access-Control-Allow-Credentials", "true")
			resspo.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if req.Method == http.MethodOptions {
				resspo.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(resspo, req)
		})
	}
}
