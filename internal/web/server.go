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
}

type BuilderOpts func(*Server)

func WithLogger(logger *slog.Logger) BuilderOpts {
	return func(s *Server) {
		s.logger = logger
	}
}

func WithDB(db *pgxpool.Pool) BuilderOpts {
	return func(s *Server) {
		s.db = db
	}
}

func WithUserService(userService user.UserService) BuilderOpts {
	return func(s *Server) {
		s.userService = userService
	}
}

func WithGroupService(groupService group.GroupService) BuilderOpts {
	return func(s *Server) {
		s.groupService = groupService
	}
}

func WithSessionService(sessionService session.SessionService) BuilderOpts {
	return func(s *Server) {
		s.sessionService = sessionService
	}
}

func WithMux(mux *http.ServeMux) BuilderOpts {
	return func(s *Server) {
		s.mux = mux
	}
}

func NewServer(db *pgxpool.Pool, opts ...BuilderOpts) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Start(port string) error {
	s.Route()
	s.logger.Info("Starting server", "addr", port)
	server := http.Server{
		Addr:    port,
		Handler: s.MiddlewareLogIP(s.mux),
	}
	return server.ListenAndServe()
}
