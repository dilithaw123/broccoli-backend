package web

import (
	"log/slog"
	"net/http"

	"github.com/dilithaw123/broccoli-backend/internal/group"
	"github.com/dilithaw123/broccoli-backend/internal/session"
	"github.com/dilithaw123/broccoli-backend/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

func WithSecretKey(secretKey string) BuilderOpts {
	return func(s *Server) {
		s.secretKey = secretKey
	}
}

func WithApiKey(key string) BuilderOpts {
	return func(s *Server) {
		s.apiKey = key
	}
}
