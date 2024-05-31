package api

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/xsadia/secred/pkg/api/handlers"
)

type Server struct {
	server *echo.Echo
}

func (s *Server) Listen(address string) {
	s.init()
	log.Panic(s.server.Start(address))
}

func (s *Server) init() {
	s.server.GET("/healthz", handlers.HealthZHandler())
	s.server.POST("/register", handlers.RegisterHandler())
}

func New() *Server {
	return &Server{
		server: echo.New(),
	}
}
