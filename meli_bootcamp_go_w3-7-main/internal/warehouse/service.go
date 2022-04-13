package warehouse

import (
	"context"
	"errors"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
)

// Errors
var (
	ErrNotFound = errors.New("warehouse not found")
)

type Service interface {
	Get(ctx context.Context, id int) (domain.Warehouse, error)
	GetAll(ctx context.Context) ([]domain.Warehouse, error)
	Save(ctx context.Context, w domain.Warehouse) (domain.Warehouse, error)
	Exists(ctx context.Context, wc string) bool
	Update(ctx context.Context, w domain.Warehouse) error
	Delete(ctx context.Context, id int) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	return s.repository.Get(ctx, id)
}

func (s *service) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	return s.repository.GetAll(ctx)
}

func (s *service) Save(ctx context.Context, w domain.Warehouse) (domain.Warehouse, error) {
	id, err := s.repository.Save(ctx, w)
	if err != nil {
		return domain.Warehouse{}, err
	}
	w.ID = id
	return w, err
}

func (s *service) Exists(ctx context.Context, wc string) bool {
	return s.repository.Exists(ctx, wc)
}

func (s *service) Update(ctx context.Context, w domain.Warehouse) error {
	return s.repository.Update(ctx, w)
}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}
