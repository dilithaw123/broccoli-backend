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
		session_id, err := strconv.ParseUint(session, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse session id", "error", err)
			http.Error(w, "missing session_id parameter", http.StatusBadRequest)
			return
		}
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
		if err != nil {
			s.logger.Error("Failed to upgrade connection", "error", err)
			return
		}
		s.logger.Info("New websocket connection", "ip", r.RemoteAddr)
		go s.pollSend(conn, session_id)
	}
}

// pollSend sends the latest submissions to the client every second
func (s *Server) pollSend(ws *websocket.Conn, session_id uint64) {
	defer ws.Close(websocket.StatusGoingAway, "it's over")
	ctx := ws.CloseRead(context.Background())
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Connection closed")
			return
		case <-time.After(time.Second):
			sub, err := s.userService.GetAllUserSubmissionsForSession(
				context.Background(),
				session_id,
			)
			if err != nil {
				s.logger.Error("Failed to get user submissions", "error", err)
				return
			}
			if err := wsjson.Write(context.Background(), ws, sub); err != nil {
				s.logger.Error("Failed to write message", "error", err)
				return
			}
		}
	}
}
