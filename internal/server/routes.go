package server

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))
	s.Get("/health", s.HealthHandler)
	api := s.Group("/api")
	accout := api.Group("/account")
	accout.Get("/getcaptcha", s.GetCaptcha)
	accout.Get("/register", s.RegisterHandler)
}
