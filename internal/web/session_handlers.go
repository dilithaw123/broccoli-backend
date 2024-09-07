package web

import (
	"encoding/json"
	"net/http"

	"github.com/dilithaw123/broccoli-backend/internal/session"
)

func (s *Server) handlePostSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type sessionRequest struct {
			GroupID uint64 `json:"groupId"`
		}
		var req sessionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sess := session.NewSession(req.GroupID)
		id, err := s.sessionService.CreateSession(r.Context(), sess)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		type response struct {
			ID uint64 `json:"id"`
		}
		if err := respondJSON(w, http.StatusCreated, response{ID: id}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
