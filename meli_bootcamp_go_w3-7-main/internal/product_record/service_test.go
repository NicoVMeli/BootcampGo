package product_record

import (
	"context"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

// Mock for Product Record repository

type productRecordRepoMock struct {
	mock.Mock
}

func (m *productRecordRepoMock) Save(ctx context.Context, productRecord domain.ProductRecord) (int, error) {
	args := m.Called(ctx, productRecord)
	return args.Int(0), args.Error(1)
}

func (m *productRecordRepoMock) Get(ctx context.Context, id int) (domain.ProductRecord, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.ProductRecord), args.Error(1)
}

// Mock for Product service

type productServiceMock struct {
	mock.Mock
}

func (m *productServiceMock) GetAll(ctx context.Context) ([]domain.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *productServiceMock) Get(ctx context.Context, id int) (domain.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *productServiceMock) Save(ctx context.Context, product domain.Product) (domain.Product, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *productServiceMock) Update(ctx context.Context, id int, product domain.Product) (domain.Product, error) {
	args := m.Called(ctx, id, product)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *productServiceMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *productServiceMock) GetRecordReportsByProductId(ctx context.Context, id int) (string, int, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(string), args.Get(1).(int), args.Error(2)
}

// Tests

func TestCreateOk(t *testing.T) {
	// ARRANGE --------------------------------
	localTime := time.Date(2022, 02, 07, 16, 57, 52, 0, time.Local) // Time is created using local system timezone
	productRecordToCreate := domain.ProductRecord{
		ID:             0,
		LastUpdateDate: &localTime,
		SalePrice:      23.1,
		PurchasePrice:  23.5,
		ProductID:      1,
	}

	expectedProductRecordServiceResult := domain.ProductRecord{
		ID:             5,
		LastUpdateDate: &localTime,
		SalePrice:      23.1,
		PurchasePrice:  23.5,
		ProductID:      1,
	}

	productRecordRepository := new(productRecordRepoMock)
	productRecordRepository.On("Save", mock.Anything, productRecordToCreate).Return(5, nil)
	productService := new(productServiceMock)
	productService.On("Get", mock.Anything, 1).Return(domain.Product{}, nil) // Returns an existing Product and no error
	productRecordService := NewService(productRecordRepository, productService)
	ctx := context.Background()

	// ACT -----------------------------
	actualServiceResult, actualServiceError := productRecordService.Save(ctx, productRecordToCreate)

	// ASSERT --------------------------
	assert.Equal(t, expectedProductRecordServiceResult, actualServiceResult)
	assert.Nil(t, actualServiceError)
}
