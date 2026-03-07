package server

import (
	"log"
	"net/http"
	"time"

	"go1fl-sprint6-final-tpl/internal/handlers"
)

type Server struct {
	logger *log.Logger
	server *http.Server
}

func NewServer(logger *log.Logger) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/upload", handlers.UploadHandler)

	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		logger: logger,
		server: httpServer,
	}
}

func (s *Server) Start() error {
	s.logger.Println("Server starting on :8080")
	return s.server.ListenAndServe()
}
