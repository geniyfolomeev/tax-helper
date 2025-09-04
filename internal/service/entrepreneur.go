package service

import (
	"errors"
	"tax-helper/internal/domain"
)

type EntrepreneurRepo interface {
	GetByID(id uint) (*domain.Entrepreneur, error)
	Create(e *domain.Entrepreneur) error
}

type TasksRepo interface {
	CreateBatch(tasks []*domain.Task) error
}

type EntrepreneurService struct {
	EntrepreneurRepo EntrepreneurRepo
	TasksRepo        TasksRepo
}

func NewEntrepreneurService(e EntrepreneurRepo, t TasksRepo) *EntrepreneurService {
	return &EntrepreneurService{EntrepreneurRepo: e, TasksRepo: t}
}

// TODO: транзакции!
func (s *EntrepreneurService) CreateEntrepreneur(e *domain.Entrepreneur) error {
	err := e.Validate()
	if err != nil {
		return err
	}
	_, err = s.EntrepreneurRepo.GetByID(e.TelegramID)
	if err != nil && !errors.Is(err, domain.ErrEntrepreneurNotFound) {
		return err
	}
	err = s.EntrepreneurRepo.Create(e)
	if err != nil {
		return err
	}
	nextDeclarationDate := e.CalculateNextDeclarationDate()
	addIncomeTask := domain.NewTask(e.TelegramID, "ready", "add_income", nextDeclarationDate)
	submitDeclarationTask := domain.NewTask(e.TelegramID, "ready", "submit declaration", nextDeclarationDate)
	return s.TasksRepo.CreateBatch([]*domain.Task{addIncomeTask, submitDeclarationTask})
}
