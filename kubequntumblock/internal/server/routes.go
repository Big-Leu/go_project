package server

import (
	"net/http"

	"kubequntumblock/controllers"
	"kubequntumblock/middleware"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.POST("/adduser",controllers.Signup)

	r.POST("/login",controllers.Login)

	r.GET("/validate",middleware.RequireAuth, controllers.Validate)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
