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

	schoolGroup := s.server.Group("/schools")
	schoolGroup.POST("", handlers.CreateSchoolHandler)
	schoolGroup.GET("", handlers.ListSchoolHandler)
	schoolGroup.POST("/:id/orders", handlers.CreateSchoolOrderHandler)

	itemGroup := s.server.Group("/items")
	itemGroup.POST("", handlers.CreateItemHandler)
	itemGroup.GET("", handlers.ListItemsHandler)

	orderGroup := s.server.Group("/orders")
	orderGroup.POST("", handlers.CreateOrderHandler)
}

func New() *Server {
	return &Server{
		server: echo.New(),
	}
}
