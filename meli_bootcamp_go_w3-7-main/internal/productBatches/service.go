package productBatches

import (
	"context"
	"errors"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
)

// Errors
var (
	ErrNotFound      = errors.New("el product batches no fue encontrada")
	ErrAlreadyExists = errors.New("el product batches  ya existe")
)

type Service interface {
	GetBySectionNumber(ctx context.Context, sectionNumber int) (domain.Section, error)
	GetQuantity(ctx context.Context, quant int) (domain.ProductBatches, error)
	Save(ctx context.Context, sec domain.ProductBatches) (domain.ProductBatches, error)
	Exists(ctx context.Context, batchNumber int) bool
}

type service struct {
	productBatchesRepository Repository
}

func NewService(sectionRepo Repository) Service {
	return &service{
		productBatchesRepository: sectionRepo,
	}
}

func (s *service) Save(ctx context.Context, sect domain.ProductBatches) (domain.ProductBatches, error) {
	exist := s.productBatchesRepository.Exists(ctx, sect.BatchNumber)
	if exist {
		return domain.ProductBatches{}, ErrAlreadyExists
	}
	newSectionId, err := s.productBatchesRepository.Save(ctx, sect)
	if err != nil {
		return domain.ProductBatches{}, err
	}
	sect.ID = newSectionId
	return sect, nil
}

func (s *service) Exists(ctx context.Context, sec int) bool {
	return s.productBatchesRepository.Exists(ctx, sec)
}

func (s *service) GetBySectionNumber(ctx context.Context, sec int) (domain.Section, error) {
	p, err := s.productBatchesRepository.GetBySectionNumber(ctx, sec)
	if err != nil {
		return domain.Section{}, ErrNotFound
	}
	return p, nil
}

func (s *service) GetQuantity(ctx context.Context, sec int) (domain.ProductBatches, error) {
	p, err := s.productBatchesRepository.GetQuantity(ctx, sec)
	if err != nil {
		return domain.ProductBatches{}, ErrNotFound
	}
	return p, nil
}
