package server

import (
	"abdi/task-manager/internal/models"
	"context"
	"net/http"
)

type Server struct {
	store      models.TaskStore
	router     *http.ServeMux
	httpServer *http.Server
}

func New(store models.TaskStore) *Server {
	s := &Server{
		store:  store,
		router: http.NewServeMux(),
	}

	s.router.HandleFunc("POST /tasks", s.createTask)
	s.router.HandleFunc("GET /tasks", s.listTasks)
	s.router.HandleFunc("PUT /tasks/{id}", s.completeTask)
	s.router.HandleFunc("DELETE /tasks/{id}", s.deleteTask)

	return s
}

func (s *Server) Start(addr string) error {
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: logging(s.router),
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
