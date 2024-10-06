package web

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/dilithaw123/broccoli-backend/internal/group"
	"github.com/dilithaw123/broccoli-backend/internal/session"
	"github.com/dilithaw123/broccoli-backend/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type room map[uint64]map[*websocket.Conn]struct{}

func newRoom() room {
	return room(make(map[uint64]map[*websocket.Conn]struct{}))
}

type sessionMap struct {
	sync.Mutex
	room
}

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
	sessions       sessionMap
}

func NewServer(db *pgxpool.Pool, opts ...BuilderOpts) *Server {
	sessions := newRoom()
	s := &Server{
		refTokenMap: make(map[string]string),
		sessions:    sessionMap{sync.Mutex{}, sessions},
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
	go s.clientUpdate()
	return server.ListenAndServe()
}
