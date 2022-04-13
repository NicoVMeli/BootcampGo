package product_record

import (
	"context"
	"errors"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/product"
)

//Errors
var (
	ErrProductNotFound = errors.New("no product matches the specified id")
)

type Service interface {
	//Get(ctx context.Context, id int) (domain.ProductRecord, error)
	Save(ctx context.Context, p domain.ProductRecord) (domain.ProductRecord, error)
}

type service struct {
	productRecordRepository Repository
	productService          product.Service
}

func NewService(productRecordRepository Repository, productService product.Service) Service {
	return &service{
		productRecordRepository: productRecordRepository,
		productService:          productService,
	}
}

func (s *service) Save(ctx context.Context, pr domain.ProductRecord) (domain.ProductRecord, error) {
	// Check if the provided product_id corresponds to an existing Product
	_, err := s.productService.Get(ctx, pr.ProductID)
	if err != nil {
		if errors.Is(err, product.ErrNotFound) {
			return domain.ProductRecord{}, ErrProductNotFound
		}
		return domain.ProductRecord{}, err
	}
	// Save the Product Record
	newProductRecordId, err := s.productRecordRepository.Save(ctx, pr)
	if err != nil {
		return domain.ProductRecord{}, err
	}
	pr.ID = newProductRecordId
	return pr, nil
}
