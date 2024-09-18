package web

import (
	"encoding/json"
	"net/http"

	"github.com/dilithaw123/broccoli-backend/internal/group"
)

func (s *Server) handleGetUserGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		groups, err := s.groupService.GetGroupsByEmail(r.Context(), email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bytes, err := json.Marshal(groups)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if _, err := w.Write(bytes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) handlePostGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Group group.Group `json:"group"`
			Email string      `json:"email"`
		}
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := s.groupService.CreateUpdateGroup(r.Context(), req.Group, req.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Server) handleAddUserToGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			UserEmail    string `json:"email"`
			GroupID      uint64 `json:"group_id"`
			RequestEmail string `json:"request_email"`
		}
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userAllowed, err := s.groupService.GroupContainsUser(
			r.Context(),
			req.GroupID,
			req.RequestEmail,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !userAllowed {
			http.Error(w, "user not allowed", http.StatusForbidden)
			return
		}

		err = s.groupService.AddUserToGroup(r.Context(), req.GroupID, req.UserEmail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
