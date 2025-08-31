package http

import (
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/http/handlers"

	"github.com/gin-gonic/gin"
)

type EntrepreneurService interface {
	CreateEntrepreneur(user *domain.Entrepreneur) error
}

type Server struct {
	engine *gin.Engine
}

func New(entrepreneurService EntrepreneurService) *Server {
	r := gin.Default()
	s := &Server{engine: r}
	s.addRoutes(entrepreneurService)
	return s
}

func (s *Server) addRoutes(entrepreneurService EntrepreneurService) {
	s.engine.POST("/register", handlers.NewEntrepreneurHandler(entrepreneurService).Register)
}

func (s *Server) Run(port string) error {
	return s.engine.Run(port)
}
