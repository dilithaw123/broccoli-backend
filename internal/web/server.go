package web

import (
	"log/slog"
	"net/http"

	"github.com/dilithaw123/broccoli-backend/internal/group"
	"github.com/dilithaw123/broccoli-backend/internal/session"
	"github.com/dilithaw123/broccoli-backend/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	db             *pgxpool.Pool
	userService    user.UserService
	groupService   group.GroupService
	sessionService session.SessionService
	mux            *http.ServeMux
	logger         *slog.Logger
	refTokenMap    map[string]string
	secretKey      string
	apiKey         string
}

func NewServer(db *pgxpool.Pool, opts ...BuilderOpts) *Server {
	s := &Server{
		refTokenMap: make(map[string]string),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Start(port string) error {
	s.Route()
	s.logger.Info("Starting server", "addr", port)
	handler := s.MiddlewareLogIP(s.mux)
	server := http.Server{
		Addr:    port,
		Handler: handler,
	}
	return server.ListenAndServe()
}
