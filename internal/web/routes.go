package web

import "net/http"

func (s *Server) Route() {
	innerMux := http.NewServeMux()
	innerMux.Handle("GET /ws/session/{id}", s.handleSessionWSConnection())
	innerMux.Handle("POST /session", s.handlePostSession())
	innerMux.Handle("POST /group/user/add", s.handleAddUserToGroup())
	innerMux.Handle("POST /group", s.handlePostGroup())
	innerMux.Handle("DELETE /group", s.handleDeleteGroup())
	innerMux.Handle("GET /user/submission", s.handleGetUserSubmission())
	innerMux.Handle("POST /user/submission", s.handlePostUserSubmission())
	innerMux.Handle("GET /user/group", s.handleGetUserGroups())
	innerMux.Handle("GET /user/authenticated", s.handleIsAuthorized())
	innerMux.Handle("GET /user", s.handleGetUser())
	innerMux.Handle("POST /user", s.handlePostUser())
	// User without access token needs to be able to hit these endpoints
	s.mux.Handle("POST /user/refresh", s.handleNewAccessToken())
	s.mux.Handle("POST /login", s.MiddlewareAPIKey(s.handleLoginSignUp()))
	s.mux.Handle("/", s.MiddlewareAuth(innerMux))
}
