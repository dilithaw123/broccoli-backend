package web

func (s *Server) Route() {
	s.mux.Handle("GET /ws/session/{id}", s.handleSessionWSConnection())
	s.mux.Handle("POST /session", s.handlePostSession())
	s.mux.Handle("POST /group/user/add", s.handleAddUserToGroup())
	s.mux.Handle("POST /group", s.handlePostGroup())
	s.mux.Handle("GET /user/submission", s.handleGetUserSubmission())
	s.mux.Handle("POST /user/submission", s.handlePostUserSubmission())
	s.mux.Handle("GET /user/group", s.handleGetUserGroups())
	s.mux.Handle("GET /user", s.handleGetUser())
	s.mux.Handle("POST /user", s.handlePostUser())
	s.mux.Handle("POST /login", s.handleLoginSignUp())
}
