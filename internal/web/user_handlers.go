package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/dilithaw123/broccoli-backend/internal/user"
)

// Get user by email or id
func (s *Server) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		id := r.URL.Query().Get("id")
		var u user.User
		var err error
		if email == "" && id == "" {
			http.Error(w, "email or id query parameter required", http.StatusBadRequest)
			return
		} else if id != "" {
			id, converr := strconv.ParseUint(id, 10, 64)
			if converr != nil {
				http.Error(w, "id query parameter must be an integer", http.StatusBadRequest)
				return
			}
			u, err = s.userService.GetUserByID(r.Context(), id)
		} else {
			u, err = s.userService.GetUserByEmail(r.Context(), strings.ToLower(email))
		}
		if err != nil {
			if err == user.ErrUserNotFound {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		uJSON, err := u.JSON()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(uJSON); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func (s *Server) handlePostUser() http.HandlerFunc {
	type request struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		u := user.NewUser(req.Name, req.Email)
		_, err := s.userService.CreateUser(r.Context(), u)
		if err != nil {
			if err == user.ErrUserAlreadyExists {
				http.Error(w, "user already exists", http.StatusConflict)
				return
			} else {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(200)
	}
}

func (s *Server) handlePostUserSubmission() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sub user.UserSubmission
		if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if err := s.userService.CreateUpdateUserSubmission(r.Context(), sub); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) handleGetUserSubmission() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session_id, err := strconv.ParseUint(r.URL.Query().Get("session_id"), 10, 64)
		if err != nil {
			http.Error(w, "missing session_id parameter", http.StatusBadRequest)
			return
		}
		var sub interface{}
		if r.URL.Query().Get("all") != "true" {
			user_id, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 64)
			if err != nil {
				http.Error(w, "missing user_id parameter", http.StatusBadRequest)
				return
			}
			sub, err = s.userService.GetUserSubmission(r.Context(), user_id, session_id)
			if err != nil {
				if err == user.ErrorUserSubmissionNotFound {
					http.Error(w, "user submission not found", http.StatusNotFound)
					return
				}
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			sub, err = s.userService.GetAllUserSubmissionsForSession(r.Context(), session_id)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(sub); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func (s *Server) handleLoginSignUp() http.HandlerFunc {
	type loginRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	type loginResponse struct {
		User         user.User `json:"user"`
		RefreshToken string    `json:"refresh_token"`
		AccessToken  string    `json:"access_token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		req.Email = strings.ToLower(req.Email)
		u, err := s.userService.GetUserByEmail(r.Context(), req.Email)
		if err != nil && !errors.Is(err, user.ErrUserNotFound) {
			s.logger.Error("Login/SignUp", "Error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if errors.Is(err, user.ErrUserNotFound) {
			u = user.NewUser(req.Name, req.Email)
			if u, err = s.userService.CreateUser(r.Context(), u); err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		refreshToken := GenerateRefreshToken()
		accessToken := GenerateAccessToken(s.secretKey)
		s.refTokenMap[u.Email] = refreshToken
		resp := loginResponse{
			User:         u,
			RefreshToken: refreshToken,
			AccessToken:  accessToken,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

// Sends ok response. If it reaches this endpoint this must mean the middleware has authenticated the user.
func (s *Server) handleIsAuthorized() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}
}

// If user's email exists in map with matching refresh token, send a new access token
func (s *Server) handleNewAccessToken() http.HandlerFunc {
	type request struct {
		Email        string `json:"email"`
		RefreshToken string `json:"refresh_token"`
	}
	type response struct {
		AccessToken string `json:"access_token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if val, ok := s.refTokenMap[req.Email]; !ok || req.RefreshToken != val {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		resp := response{AccessToken: GenerateAccessToken(s.secretKey)}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
