package web

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func (s *Server) handleSessionWSConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := r.PathValue("id")
		sessionId, err := strconv.ParseUint(session, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse session id", "error", err)
			http.Error(w, "missing session_id parameter", http.StatusBadRequest)
			return
		}
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
		if err != nil {
			s.logger.Error("Failed to upgrade connection", "error", err)
			return
		}
		s.logger.Info("New websocket connection", "ip", r.RemoteAddr)
		s.addToSessionMap(sessionId, conn)
	}
}

func (s *Server) clientUpdate() {
	for {
		time.After(time.Second)
		func() {
			s.sessions.Lock()
			defer s.sessions.Unlock()
			for k, v := range s.sessions.room {
				ctx := context.Background()
				ctx = context.WithValue(ctx, "sessionId", k)
				sub, err := s.userService.GetAllUserSubmissionsForSession(ctx, k)
				if err != nil {
					s.logger.Error("Failed to get user submissions", "error", err, "sessionId", k)
					continue
				}
				for conn := range v {
					if err := wsjson.Write(ctx, conn, sub); err != nil {
						s.logger.Error("Failed to write message", "error", err)
						delete(v, conn)
						continue
					}
				}
				if len(v) == 0 {
					delete(s.sessions.room, k)
				}
			}
		}()
	}
}

func (s *Server) addToSessionMap(sessionId uint64, ws *websocket.Conn) {
	s.sessions.Lock()
	defer s.sessions.Unlock()
	if room, ok := s.sessions.room[sessionId]; ok {
		room[ws] = struct{}{}
	} else {
		room := make(map[*websocket.Conn]struct{})
		room[ws] = struct{}{}
		s.sessions.room[sessionId] = room
	}
}
