package product

import (
	"context"
	"fmt"
	"testing"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type dbMock struct {
	mock.Mock
}

func (m *dbMock) GetAll(ctx context.Context) ([]domain.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *dbMock) Get(ctx context.Context, id int) (domain.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *dbMock) Exists(ctx context.Context, productCode string) bool {
	args := m.Called(ctx, productCode)
	return args.Bool(0)
}

func (m *dbMock) Save(ctx context.Context, product domain.Product) (int, error) {
	args := m.Called(ctx, product)
	return args.Int(0), args.Error(1)
}

func (m *dbMock) Update(ctx context.Context, product domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *dbMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *dbMock) GetRecordReportsByProductId(ctx context.Context, productId int) (recordsCount int, err error) {
	args := m.Called(ctx, productId)
	return args.Int(0), args.Error(1)
}

func TestCreateOk(t *testing.T) {
	// ARRANGE
	productToCreate := domain.Product{
		ID:             0,
		Description:    "Yogurt",
		ExpirationRate: 1,
		FreezingRate:   2,
		Height:         6.4,
		Length:         4.5,
		Netweight:      3.4,
		ProductCode:    "PROD08",
		RecomFreezTemp: 1.3,
		Width:          1.2,
		ProductTypeID:  2,
		SellerID:       4,
	}

	expectedServiceResult := domain.Product{
		ID:             1,
		Description:    "Yogurt",
		ExpirationRate: 1,
		FreezingRate:   2,
		Height:         6.4,
		Length:         4.5,
		Netweight:      3.4,
		ProductCode:    "PROD08",
		RecomFreezTemp: 1.3,
		Width:          1.2,
		ProductTypeID:  2,
		SellerID:       4,
	}

	productRepository := new(dbMock)
	productRepository.On("Exists", mock.Anything, productToCreate.ProductCode).Return(false)
	productRepository.On("Save", mock.Anything, productToCreate).Return(1, nil)
	productService := NewService(productRepository)
	ctx := context.Background()

	// ACT
	actualServiceResult, actualServiceError := productService.Save(ctx, productToCreate)

	// ASSERT
	assert.Equal(t, expectedServiceResult, actualServiceResult)
	assert.Nil(t, actualServiceError)
}

func TestCreateConflict(t *testing.T) {
	// ARRANGE
	productToCreate := domain.Product{
		ID:             0,
		Description:    "Yogurt",
		ExpirationRate: 1,
		FreezingRate:   2,
		Height:         6.4,
		Length:         4.5,
		Netweight:      3.4,
		ProductCode:    "PROD08",
		RecomFreezTemp: 1.3,
		Width:          1.2,
		ProductTypeID:  2,
		SellerID:       4,
	}

	expectedServiceResult := domain.Product{}
	expectedServiceError := fmt.Errorf("product already exists")
	productRepository := new(dbMock)
	productRepository.On("Exists", mock.Anything, productToCreate.ProductCode).Return(true)
	productService := NewService(productRepository)
	ctx := context.Background()

	// ACT
	actualServiceResult, actualServiceError := productService.Save(ctx, productToCreate)

	//ASSERT
	assert.Equal(t, expectedServiceResult, actualServiceResult)
	assert.Equal(t, expectedServiceError, actualServiceError)
}

func TestFindAll(t *testing.T) {
	// ARRANGE
	mockedRepoResult := []domain.Product{}
	mockedRepoResult = append(mockedRepoResult, domain.Product{
		ID:             1,
		Description:    "Yogurt",
		ExpirationRate: 1,
		FreezingRate:   2,
		Height:         6.4,
		Length:         4.5,
		Netweight:      3.4,
		ProductCode:    "PROD08",
		RecomFreezTemp: 1.3,
		Width:          1.2,
		ProductTypeID:  2,
		SellerID:       4,
	})

	expectedServiceResult := []domain.Product{}
	expectedServiceResult = append(expectedServiceResult, domain.Product{
		ID:             1,
		Description:    "Yogurt",
		ExpirationRate: 1,
		FreezingRate:   2,
		Height:         6.4,
		Length:         4.5,
		Netweight:      3.4,
		ProductCode:    "PROD08",
		RecomFreezTemp: 1.3,
		Width:          1.2,
		ProductTypeID:  2,
		SellerID:       4,
	})
	productRepository := new(dbMock)
	productRepository.On("GetAll", mock.Anything).Return(mockedRepoResult, nil)
	productService := NewService(productRepository)
	ctx := context.Background()

	// ACT
	actualServiceResult, actualServiceError := productService.GetAll(ctx)

	// ASSERT
	assert.True(t, len(actualServiceResult) == 1)
	assert.Equal(t, expectedServiceResult, actualServiceResult)
	assert.Nil(t, actualServiceError)
}

func TestFindByIdNonExistent(t *testing.T) {
	// ARRANGE
	mockedRepoResult := domain.Product{}
	mockedRepoError := fmt.Errorf("sql: no rows in result set")

	expectedServiceResult := domain.Product{}
	expectedServiceError := fmt.Errorf("product not found")

	productRepository := new(dbMock)
	productRepository.On("Get", mock.Anything, 1).Return(mockedRepoResult, mockedRepoError)
	productService := NewService(productRepository)
	ctx := context.Background()

	// ACT
	actualServiceResult, actualServiceError := productService.Get(ctx, 1)

	// ASSERT
	assert.Equal(t, expectedServiceResult, actualServiceResult)
	assert.Equal(t, expectedServiceError, actualServiceError)
}

func TestFindByIdExistent(t *testing.T) {
	// ARRANGE
	mockedRepoResult := domain.Product{
		ID:             1,
		Description:    "Yogurt",
		ExpirationRate: 1,
		FreezingRate:   2,
		Height:         6.4,
		Length:         4.5,
		Netweight:      3.4,
		ProductCode:    "PROD08",
		RecomFreezTemp: 1.3,
		Width:          1.2,
		ProductTypeID:  2,
		SellerID:       4,
	}

	expectedServiceResult := domain.Product{
		ID:             1,
		Description:    "Yogurt",
		ExpirationRate: 1,
		FreezingRate:   2,
		Height:         6.4,
		Length:         4.5,
		Netweight:      3.4,
		ProductCode:    "PROD08",
		RecomFreezTemp: 1.3,
		Width:          1.2,
		ProductTypeID:  2,
		SellerID:       4,
	}

	productRepository := new(dbMock)
	productRepository.On("Get", mock.Anything, 1).Return(mockedRepoResult, nil)
	productService := NewService(productRepository)
	ctx := context.Background()

	// ACT
	actualServiceResult, actualServiceError := productService.Get(ctx, 1)

	// ASSERT
	assert.Equal(t, expectedServiceResult, actualServiceResult)
	assert.Nil(t, actualServiceError)
}

func TestUpdateExistent(t *testing.T) {
	// ARRANGE
	mockedRepoResultGet := domain.Product{
		ID:             1,
		Description:    "Yogurt",
		ExpirationRate: 1,
		FreezingRate:   2,
		Height:         6.4,
		Length:         4.5,
		Netweight:      3.4,
		ProductCode:    "PROD08",
		RecomFreezTemp: 1.3,
		Width:          1.2,
		ProductTypeID:  2,
		SellerID:       4,
	}

	productPatchToApplyId := 1
	productPatchToApply := domain.Product{
		Description:    "New Yogurt",
		ExpirationRate: 3,
		FreezingRate:   4,
		Height:         6.5,
		Length:         4.6,
		Netweight:      3.8,
		ProductCode:    "NEWPROD08",
		RecomFreezTemp: 1.1,
		Width:          1.5,
		ProductTypeID:  3,
		SellerID:       6,
	}

	patchedProductToUpdate := domain.Product{
		ID:             1,
		Description:    "New Yogurt",
		ExpirationRate: 3,
		FreezingRate:   4,
		Height:         6.5,
		Length:         4.6,
		Netweight:      3.8,
		ProductCode:    "NEWPROD08",
		RecomFreezTemp: 1.1,
		Width:          1.5,
		ProductTypeID:  3,
		SellerID:       6,
	}
	expectedServiceResult := patchedProductToUpdate

	productRepository := new(dbMock)
	productRepository.On("Get", mock.Anything, productPatchToApplyId).Return(mockedRepoResultGet, nil)
	productRepository.On("Update", mock.Anything, patchedProductToUpdate).Return(nil)
	productService := NewService(productRepository)
	ctx := context.Background()

	// ACT
	actualServiceResult, actualServiceError := productService.Update(ctx, productPatchToApplyId, productPatchToApply)

	// ASSERT
	assert.Equal(t, expectedServiceResult, actualServiceResult)
	assert.Nil(t, actualServiceError)
}

func TestUpdateNonExistent(t *testing.T) {
	// ARRANGE
	mockedRepoResult := domain.Product{}
	mockedRepoError := fmt.Errorf("sql: no rows in result set")

	productPatchToApplyId := 1
	productPatchToApply := domain.Product{
		Description:    "New Yogurt",
		ExpirationRate: 3,
		FreezingRate:   4,
		ProductCode:    "NEWPROD08",
	}

	expectedServiceResult := domain.Product{}
	expectedServiceError := fmt.Errorf("product not found")

	productRepository := new(dbMock)
	productRepository.On("Get", mock.Anything, productPatchToApplyId).Return(mockedRepoResult, mockedRepoError)
	productService := NewService(productRepository)
	ctx := context.Background()

	// ACT
	actualServiceResult, actualServiceError := productService.Update(ctx, productPatchToApplyId, productPatchToApply)

	// ASSERT
	assert.Equal(t, expectedServiceResult, actualServiceResult)
	assert.Equal(t, expectedServiceError, actualServiceError)
}

func TestDeleteNonExistent(t *testing.T) {
	// ARRANGE
	mockedRepoError := fmt.Errorf("product not found")
	productToDeleteId := 1
	expectedServiceError := fmt.Errorf("product not found")

	productRepository := new(dbMock)
	productRepository.On("Delete", mock.Anything, productToDeleteId).Return(mockedRepoError)
	productService := NewService(productRepository)
	ctx := context.Background()

	// ACT
	actualServiceError := productService.Delete(ctx, productToDeleteId)

	// ASSERT
	assert.Equal(t, expectedServiceError, actualServiceError)
}

func TestDeleteOK(t *testing.T) {
	// ARRANGE
	productToDeleteId := 1

	productRepository := new(dbMock)
	productRepository.On("Delete", mock.Anything, productToDeleteId).Return(nil)
	productService := NewService(productRepository)
	ctx := context.Background()

	// ACT
	actualServiceError := productService.Delete(ctx, productToDeleteId)

	// ASSERT
	assert.Nil(t, actualServiceError)
}
