package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"tax-helper/config"
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/http/handlers"

	"github.com/gin-gonic/gin"
)

type ServerWrapper struct {
	server *http.Server
}

type EntrepreneurService interface {
	CreateEntrepreneur(user *domain.Entrepreneur) error
}

func cancelMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-ctx.Done():
			fmt.Println("connection closed by peer")
		case <-done:
		}
	}
}

func NewServer(cfg *config.Config, entrepreneurService EntrepreneurService) *ServerWrapper {
	engine := gin.Default()
	engine.Use(cancelMiddleware())
	engine.POST("/register", handlers.NewEntrepreneurHandler(entrepreneurService).Register)
	s := &http.Server{
		Handler: engine,
		Addr:    cfg.HTTPPort,
	}
	return &ServerWrapper{server: s}
}

func (s *ServerWrapper) Run(ctx context.Context) error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	<-ctx.Done()
	return s.server.Shutdown(ctx)
}
