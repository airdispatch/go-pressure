package pressure

import (
	"net/http"
)

type Server struct {
	backend *http.Server
	router  *Router
	Debug   bool
	Port    string
	*Logger
}

// ---- SERVER METHODS ---- //

func CreateServer(port string, debug bool) *Server {
	s := &Server{}
	s.Port = port

	s.Logger = NewLogger(ERROR)

	s.Debug = debug
	if s.Debug {
		s.Logger.LogLevel = DEBUG
	}
	s.router = &Router{}

	return s
}

func (s *Server) RunServer() {
	s.backend = &http.Server{}
	s.backend.Handler = s.router
	s.backend.Addr = s.Port

	s.router.Logger = s.Logger

	s.Logger.LogDebug("Server is now running at", s.Port)
	err := s.backend.ListenAndServe()
	if err != nil {
		s.Logger.LogError("Error Occured to Run Server:", err)
		return
	}
}

func (s *Server) RegisterURL(url_tuple ...Route) {
	if s.router == nil {
	}
	for _, u := range url_tuple {
		s.router.AddRoute(u)
	}
}
