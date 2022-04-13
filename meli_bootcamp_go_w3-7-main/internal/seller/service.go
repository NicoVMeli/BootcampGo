package seller

import (
	"context"
	"errors"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
)

// Errors
var (
	ErrNotFound = errors.New("seller not found")
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Seller, error)
	Get(ctx context.Context, id int) (domain.Seller, error)
	Save(ctx context.Context, se domain.Seller) (domain.Seller, error)
	Update(ctx context.Context, se domain.Seller) (domain.Seller, error)
	Delete(ctx context.Context, id int) error
	Exists(ctx context.Context, cid int) bool
}

type service struct {
	repo Repository
}

func NewService(s Repository) Service {
	return &service{repo: s}
}

func (s *service) GetAll(ctx context.Context) ([]domain.Seller, error) {
	ps, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

func (s *service) Get(ctx context.Context, id int) (domain.Seller, error) {
	p, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Seller{}, err
	}

	return p, nil
}

func (s *service) Save(ctx context.Context, se domain.Seller) (domain.Seller, error) {
	exist := s.repo.Exists(ctx, se.CID)
	if exist {
		return domain.Seller{}, errors.New("el seller ya existe")
	}
	p, err := s.repo.Save(ctx, se)

	if err != nil {
		return domain.Seller{}, err
	}

	se.ID = p

	return se, nil
}

func (s *service) Update(ctx context.Context, se domain.Seller) (domain.Seller, error) {

	err := s.repo.Update(ctx, se)

	if err != nil {
		return domain.Seller{}, err
	}

	return se, nil

}

func (s *service) Delete(ctx context.Context, id int) error {

	return s.repo.Delete(ctx, id)

}

func (s *service) Exists(ctx context.Context, sec int) bool {
	return s.repo.Exists(ctx, sec)
}
