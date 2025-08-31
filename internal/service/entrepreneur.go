package service

import (
	"errors"
	"tax-helper/internal/domain"
)

type EntrepreneurRepo interface {
	GetByID(id uint) (*domain.Entrepreneur, error)
	Create(e *domain.Entrepreneur) error
}

type EntrepreneurService struct {
	repo EntrepreneurRepo
}

func NewEntrepreneurService(r EntrepreneurRepo) *EntrepreneurService {
	return &EntrepreneurService{repo: r}
}

func (s *EntrepreneurService) CreateEntrepreneur(e *domain.Entrepreneur) error {
	err := e.Validate()
	if err != nil {
		return err
	}
	_, err = s.repo.GetByID(e.TelegramID)
	if err != nil && !errors.Is(err, domain.ErrEntrepreneurNotFound) {
		return err
	}
	return s.repo.Create(e)
}
