package api

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/xsadia/secred/pkg/api/handlers"
	"github.com/xsadia/secred/pkg/middlewares"
)

type Server struct {
	server *echo.Echo
}

func (s *Server) Listen(address string) {
	s.init()
	log.Panic(s.server.Start(address))
}

func (s *Server) init() {
	s.server.Use(middlewares.AuthMiddleware)
	s.server.GET("/healthz", handlers.HealthZHandler)
	s.server.POST("/register", handlers.RegisterHandler)
	s.server.POST("/login", handlers.LoginHandler)
	s.server.POST("/schools", handlers.CreateSchoolHandler)
	s.server.POST("/items", handlers.CreateItemHandler)
	s.server.GET("/items", handlers.ListItemsHandler)
}

func New() *Server {
	return &Server{
		server: echo.New(),
	}
}
