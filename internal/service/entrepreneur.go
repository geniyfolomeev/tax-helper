package service

import (
	"errors"
	"tax-helper/internal/domain"
)

type EntrepreneurRepo interface {
	GetByID(id uint) (*domain.Entrepreneur, error)
}

type EntrepreneurTasksRepo interface {
	Create(e *domain.Entrepreneur, tasks []*domain.Task) error
}

type EntrepreneurService struct {
	EntrepreneurRepo      EntrepreneurRepo
	EntrepreneurTasksRepo EntrepreneurTasksRepo
}

func NewEntrepreneurService(e EntrepreneurRepo, et EntrepreneurTasksRepo) *EntrepreneurService {
	return &EntrepreneurService{EntrepreneurRepo: e, EntrepreneurTasksRepo: et}
}

func (s *EntrepreneurService) CreateEntrepreneur(e *domain.Entrepreneur) error {
	err := e.Validate()
	if err != nil {
		return err
	}
	_, err = s.EntrepreneurRepo.GetByID(e.TelegramID)
	if err != nil && !errors.Is(err, domain.ErrEntrepreneurNotFound) {
		return err
	}
	nextDeclarationDate := e.CalculateNextDeclarationDate()
	addIncomeTask := domain.NewTask(e.TelegramID, "ready", "add_income", nextDeclarationDate)
	submitDeclarationTask := domain.NewTask(e.TelegramID, "ready", "submit declaration", nextDeclarationDate)
	return s.EntrepreneurTasksRepo.Create(e, []*domain.Task{addIncomeTask, submitDeclarationTask})
}
