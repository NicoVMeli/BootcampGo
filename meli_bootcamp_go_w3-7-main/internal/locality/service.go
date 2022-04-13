package locality

import (
	"context"
	"errors"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
)

var (
	ErrNotFound = errors.New("locality not found")
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Locality, error)
	GetByZipCode(ctx context.Context, zipCode string) (domain.Locality, error)
	GetSellers(ctx context.Context, l domain.Locality) ([]domain.Seller, error)
	Save(ctx context.Context, lo domain.Locality) (domain.Locality, error)
	Exists(ctx context.Context, zipC string) bool
}
type service struct {
	repo Repository
}

func NewService(l Repository) Service {
	return &service{repo: l}
}
func (l *service) GetAll(ctx context.Context) ([]domain.Locality, error) {
	ls, err := l.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return ls, nil
}
func (l *service) Save(ctx context.Context, lo domain.Locality) (domain.Locality, error) {
	exist := l.repo.Exists(ctx, lo.ZipCode)
	if exist {
		return domain.Locality{}, ErrNotFound
	}
	p, err := l.repo.Save(ctx, lo)
	if err != nil {
		return domain.Locality{}, err
	}
	lo.ID = p
	return lo, nil
}
func (l *service) Exists(ctx context.Context, zipC string) bool {
	return l.repo.Exists(ctx, zipC)
}
func (l *service) GetByZipCode(ctx context.Context, zipCode string) (domain.Locality, error) {
	//return l.repo.GetByZipCode(ctx, zipCode)
	p, err := l.repo.GetByZipCode(ctx, zipCode)
	if err != nil {
		return domain.Locality{}, err
	}
	return p, nil
}
func (s *service) GetSellers(ctx context.Context, l domain.Locality) ([]domain.Seller, error) {
	return s.repo.GetSellers(ctx, l)
}
