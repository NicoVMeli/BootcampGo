package carry

import (
	"context"
	"errors"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
)

// Errors
var (
	ErrNotFound = errors.New("carry not found")
)

type Service interface {
	Get(ctx context.Context, id int) (domain.Carry, error)
	GetAll(ctx context.Context) ([]domain.Carry, error)
	Save(ctx context.Context, c domain.Carry) (domain.Carry, error)
	Exists(ctx context.Context, cc string) bool
	Update(ctx context.Context, c domain.Carry) error
	Delete(ctx context.Context, id int) error
	GetCarriesReportByLocality(ctx context.Context) ([]domain.CarriesReport, error)
	GetCarryReportByLocalityId(ctx context.Context, id int) (domain.CarriesReport, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Get(ctx context.Context, id int) (domain.Carry, error) {
	return s.repository.Get(ctx, id)
}

func (s *service) GetAll(ctx context.Context) ([]domain.Carry, error) {
	return s.repository.GetAll(ctx)
}

func (s *service) Save(ctx context.Context, c domain.Carry) (domain.Carry, error) {
	id, err := s.repository.Save(ctx, c)
	if err != nil {
		return domain.Carry{}, err
	}
	c.ID = id
	return c, err
}

func (s *service) Exists(ctx context.Context, cc string) bool {
	return s.repository.Exists(ctx, cc)
}

func (s *service) Update(ctx context.Context, c domain.Carry) error {
	return s.repository.Update(ctx, c)
}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}

func (s *service) GetCarriesReportByLocality(ctx context.Context) ([]domain.CarriesReport, error) {
	return s.repository.GetCarriesByLocality(ctx)
}

func (s *service) GetCarryReportByLocalityId(ctx context.Context, id int) (domain.CarriesReport, error) {
	return s.repository.GetCarriesByLocalityId(ctx, id)
}
