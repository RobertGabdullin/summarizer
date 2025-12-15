package server

import "net/http"

type Server struct {
	server http.Server
}

func New(addr string, handlers []http.Handler, patterns []string) *Server {
	mux := http.NewServeMux()

	for i := range len(handlers) {
		mux.Handle(patterns[i], handlers[i])
	}

	return &Server{
		server: http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}
