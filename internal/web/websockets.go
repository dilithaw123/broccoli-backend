package web

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type userChange struct {
	UserId uint64 `json:"user_id"`
}

func (s *Server) handleSessionWSConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := r.PathValue("id")
		sessionId, err := strconv.ParseUint(session, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse session id", "error", err)
			http.Error(w, "missing session_id parameter", http.StatusBadRequest)
			return
		}
		email := r.Context().Value("email").(string)
		exists, err := s.sessionService.UserInSession(r.Context(), sessionId, email)
		if err != nil {
			s.logger.Error("Failed to check if user is in session", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if !exists {
			s.logger.Info("User not in session", "sessionId", sessionId, "email", email)
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		conn, err := websocket.Accept(
			w,
			r,
			&websocket.AcceptOptions{OriginPatterns: []string{"broccoli.buzz"}},
		)
		if err != nil {
			s.logger.Error("Failed to upgrade connection", "error", err)
			return
		}
		s.logger.Info("New websocket connection", "ip", r.RemoteAddr)
		s.addToSessionMap(sessionId, conn)
		go s.readConn(context.Background(), sessionId, conn)
	}
}

func (s *Server) sendUserChange(ctx context.Context, sessionId uint64, userId uint64) error {
	s.sessions.Lock()
	defer s.sessions.Unlock()
	for conn := range s.sessions.room[sessionId] {
		if err := wsjson.Write(ctx, conn, userChange{UserId: userId}); err != nil {
			continue
		}
	}
	return nil
}

func (s *Server) readConn(ctx context.Context, sessionId uint64, conn *websocket.Conn) {
	for {
		_, bytes, err := conn.Read(ctx)
		if err != nil {
			s.logger.Info("Closing websocket connection", "error", err)
			conn.Close(websocket.StatusNormalClosure, "bye")
			s.removeFromSessionMap(sessionId, conn)
			return
		}
		var v userChange
		if err := json.Unmarshal(bytes, &v); err == nil {
			s.logger.Debug("User change", "sessionId", sessionId, "userId", v.UserId)
			s.sendUserChange(ctx, sessionId, v.UserId)
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
