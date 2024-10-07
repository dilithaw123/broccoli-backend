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
		go s.readConn(context.Background(), sessionId, conn)
	}
}

func (s *Server) readConn(ctx context.Context, sessionId uint64, conn *websocket.Conn) {
	for {
		_, _, err := conn.Read(ctx)
		if err != nil {
			s.logger.Info("Closing websocket connection", "error", err)
			conn.Close(websocket.StatusNormalClosure, "bye")
			s.removeFromSessionMap(sessionId, conn)
			return
		}
	}
}

func (s *Server) clientUpdate() {
	for {
		time.Sleep(time.Second)
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

func (s *Server) addToSessionMap(sessionId uint64, conn *websocket.Conn) {
	s.sessions.Lock()
	defer s.sessions.Unlock()
	if room, ok := s.sessions.room[sessionId]; ok {
		room[conn] = struct{}{}
	} else {
		room := make(map[*websocket.Conn]struct{})
		room[conn] = struct{}{}
		s.sessions.room[sessionId] = room
	}
}

func (s *Server) removeFromSessionMap(sessionId uint64, conn *websocket.Conn) {
	s.sessions.Lock()
	defer s.sessions.Unlock()
	delete(s.sessions.room[sessionId], conn)
}
