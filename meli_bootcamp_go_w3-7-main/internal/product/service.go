package product

import (
	"context"
	"errors"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
)

// Errors
var (
	ErrNotFound               = errors.New("product not found")
	ErrAlreadyExists          = errors.New("product already exists")
	ErrProductRecordsNotFound = errors.New("product records not found for the provided product id")
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	Get(ctx context.Context, id int) (domain.Product, error)
	Save(ctx context.Context, p domain.Product) (domain.Product, error)
	Update(ctx context.Context, id int, p domain.Product) (domain.Product, error)
	Delete(ctx context.Context, id int) error
	GetRecordReportsByProductId(ctx context.Context, productId int) (string, int, error)
}

type service struct {
	productRepository Repository
}

func NewService(productRepository Repository) Service {
	return &service{
		productRepository: productRepository,
	}
}

func (s *service) GetAll(ctx context.Context) ([]domain.Product, error) {
	products, err := s.productRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *service) Get(ctx context.Context, id int) (domain.Product, error) {
	product, err := s.productRepository.Get(ctx, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return domain.Product{}, ErrNotFound
		}
		return domain.Product{}, err
	}
	return product, nil
}

func (s *service) Save(ctx context.Context, p domain.Product) (domain.Product, error) {
	if s.productRepository.Exists(ctx, p.ProductCode) {
		return domain.Product{}, ErrAlreadyExists
	}
	newProductId, err := s.productRepository.Save(ctx, p)
	if err != nil {
		return domain.Product{}, err
	}
	p.ID = newProductId
	return p, nil
}

func (s *service) Update(ctx context.Context, id int, productPatch domain.Product) (domain.Product, error) {

	productToUpdate, err := s.productRepository.Get(ctx, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return domain.Product{}, ErrNotFound
		}
		return domain.Product{}, err
	}
	updatedProduct := updateFields(productToUpdate, productPatch)
	err = s.productRepository.Update(ctx, updatedProduct)
	if err != nil {
		return domain.Product{}, err
	}
	return updatedProduct, nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	if err := s.productRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *service) GetRecordReportsByProductId(ctx context.Context, productId int) (string, int, error) {
	product, err := s.productRepository.Get(ctx, productId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "", 0, ErrNotFound
		}
		return "", 0, err
	}
	productRecordsCount, err := s.productRepository.GetRecordReportsByProductId(ctx, productId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "", 0, ErrProductRecordsNotFound
		}
		return "", 0, err
	}
	return product.Description, productRecordsCount, nil
}

func updateFields(productToUpdate domain.Product, productPatch domain.Product) domain.Product {
	if productPatch.Description != productToUpdate.Description && productPatch.Description != "" {
		productToUpdate.Description = productPatch.Description
	}
	if productPatch.ExpirationRate != productToUpdate.ExpirationRate && productPatch.ExpirationRate != 0 {
		productToUpdate.ExpirationRate = productPatch.ExpirationRate
	}
	if productPatch.FreezingRate != productToUpdate.FreezingRate && productPatch.FreezingRate != 0 {
		productToUpdate.FreezingRate = productPatch.FreezingRate
	}
	if productPatch.Height != productToUpdate.Height && productPatch.Height != 0 {
		productToUpdate.Height = productPatch.Height
	}
	if productPatch.Length != productToUpdate.Length && productPatch.Length != 0 {
		productToUpdate.Length = productPatch.Length
	}
	if productPatch.Netweight != productToUpdate.Netweight && productPatch.Netweight != 0 {
		productToUpdate.Netweight = productPatch.Netweight
	}
	if productPatch.ProductCode != productToUpdate.ProductCode && productPatch.ProductCode != "" {
		productToUpdate.ProductCode = productPatch.ProductCode
	}
	if productPatch.RecomFreezTemp != productToUpdate.RecomFreezTemp && productPatch.RecomFreezTemp != 0 {
		productToUpdate.RecomFreezTemp = productPatch.RecomFreezTemp
	}
	if productPatch.Width != productToUpdate.Width && productPatch.Width != 0 {
		productToUpdate.Width = productPatch.Width
	}
	if productPatch.ProductTypeID != productToUpdate.ProductTypeID && productPatch.ProductTypeID != 0 {
		productToUpdate.ProductTypeID = productPatch.ProductTypeID
	}
	if productPatch.SellerID != productToUpdate.SellerID && productPatch.SellerID != 0 {
		productToUpdate.SellerID = productPatch.SellerID
	}
	return productToUpdate
}
