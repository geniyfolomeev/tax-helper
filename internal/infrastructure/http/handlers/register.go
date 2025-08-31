package handlers

import (
	"errors"
	"net/http"
	"tax-helper/internal/domain"
	"time"

	"github.com/gin-gonic/gin"
)

type EntrepreneurService interface {
	CreateEntrepreneur(user *domain.Entrepreneur) error
}

type RegisterEntrepreneurRequest struct {
	TelegramID       uint      `json:"telegram_id" binding:"required"`
	RegistrationDate time.Time `json:"registration_date" binding:"required"`
	YearTotalAmount  float64   `json:"year_total_amount" binding:"required,gte=0"`
}

type EntrepreneurHandler struct {
	service EntrepreneurService
}

func NewEntrepreneurHandler(service EntrepreneurService) *EntrepreneurHandler {
	return &EntrepreneurHandler{service: service}
}

func (h *EntrepreneurHandler) Register(c *gin.Context) {
	var req RegisterEntrepreneurRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateEntrepreneur(&domain.Entrepreneur{
		TelegramID:       req.TelegramID,
		RegistrationDate: req.RegistrationDate,
		YearTotalAmount:  req.YearTotalAmount,
		Status:           "active",
	})
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrEntrepreneurAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
