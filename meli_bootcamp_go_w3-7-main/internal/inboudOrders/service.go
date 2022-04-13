package inboudOrders

import (
	"context"
	"errors"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
)

var ErrorEmployeesNoExist = errors.New("El empleado no existe")

type Service interface {
	Save(ctx context.Context, i domain.InboudOrders) (domain.InboudOrders, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Save(ctx context.Context, i domain.InboudOrders) (domain.InboudOrders, error) {
	exists := s.repository.ExistsEmployee(ctx, i.EmployeeId)
	if !exists {
		return domain.InboudOrders{}, ErrorEmployeesNoExist
	}
	id, err := s.repository.Save(ctx, i)
	if err != nil {
		return domain.InboudOrders{}, err
	}
	i.ID = id
	return i, nil
}
