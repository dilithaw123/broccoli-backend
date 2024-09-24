package web

import (
	"net/http"
)

func (s *Server) MiddlewareLogIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("New request ", "method", r.Method, "path", r.URL.Path, "ip", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authcookie, err := r.Cookie("access_token")
		if err != nil || authcookie == nil {
			s.logger.Debug("Error getting access token", "Error", err)
			s.logger.Info("Unauthorized request", "ip", r.RemoteAddr)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !ParseAndValidateToken(authcookie.Value, s.secretKey) {
			s.logger.Info(
				"Unauthorized access token",
				"ip",
				r.RemoteAddr,
				"token",
				authcookie.Value,
			)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) MiddlewareAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, ok := r.Header["X-Api-Key"]
		if !ok || len(apiKey) == 0 || apiKey[0] != s.apiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
