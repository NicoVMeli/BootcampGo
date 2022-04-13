package warehouse

import (
	"context"
	"errors"
	"testing"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repoM struct {
	mock.Mock
}

func (m *repoM) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Warehouse), args.Error(1)
}

func (m *repoM) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Warehouse), args.Error(1)
}

func (m *repoM) Exists(ctx context.Context, warehouseCode string) bool {
	args := m.Called(ctx, warehouseCode)
	return args.Bool(0)
}

func (m *repoM) Save(ctx context.Context, w domain.Warehouse) (int, error) {
	args := m.Called(ctx, w)
	return args.Int(0), args.Error(1)
}

func (m *repoM) Update(ctx context.Context, w domain.Warehouse) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *repoM) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateWarehouseOK(t *testing.T) {

	repo := new(repoM)
	repo.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	repo.On("Exists", mock.Anything, mock.Anything).Return(false)
	expectedResult := domain.Warehouse{
		ID:                 1,
		Address:            "Monroe 860",
		Telephone:          "47470000",
		WarehouseCode:      "CAD",
		MinimumCapacity:    50,
		MinimumTemperature: 10,
	}
	serviceT := NewService(repo)
	ctx := context.Background()
	warehouseRequest := domain.Warehouse{
		Address:            "Monroe 860",
		Telephone:          "47470000",
		WarehouseCode:      "CAD",
		MinimumCapacity:    50,
		MinimumTemperature: 10,
	}
	result, _ := serviceT.Save(ctx, warehouseRequest)
	assert.Equal(t, result, expectedResult)
}

func TestCreateWarehouseConflict(t *testing.T) {
	repo := new(repoM)
	repo.On("Exists", mock.Anything, mock.Anything).Return(true)
	serviceT := NewService(repo)
	ctx := context.Background()
	result := serviceT.Exists(ctx, domain.Warehouse{}.WarehouseCode)
	assert.True(t, result)
}

func TestFindAllWarehouses(t *testing.T) {
	expectedResult := []domain.Warehouse{}
	expectedResult = append(expectedResult,
		domain.Warehouse{
			ID:                 1,
			Address:            "Monroe 860",
			Telephone:          "47470000",
			WarehouseCode:      "CAD",
			MinimumCapacity:    50,
			MinimumTemperature: 10,
		},
		domain.Warehouse{
			ID:                 2,
			Address:            "Casa 20",
			Telephone:          "7223828",
			WarehouseCode:      "SHA",
			MinimumCapacity:    10,
			MinimumTemperature: 7,
		},
	)
	repo := new(repoM)
	repo.On("GetAll", mock.Anything).Return(expectedResult, nil)
	serviceT := NewService(repo)
	ctx := context.Background()
	results, err := serviceT.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestFindWarehouseByNonExistentId(t *testing.T) {
	expectedMessageError := "sql: no rows in result set"
	repo := new(repoM)
	repo.On("Get", mock.Anything, mock.Anything).Return(domain.Warehouse{}, errors.New(expectedMessageError))
	serviceT := NewService(repo)
	ctx := context.Background()
	result, err := serviceT.Get(ctx, 2)
	assert.Equal(t, expectedMessageError, err.Error())
	assert.Empty(t, result)
}

func TestFindWarehouseByExistentId(t *testing.T) {
	expectedResult := domain.Warehouse{
		ID:                 2,
		Address:            "Casa 20",
		Telephone:          "7223828",
		WarehouseCode:      "SHA",
		MinimumCapacity:    10,
		MinimumTemperature: 7,
	}
	repo := new(repoM)
	repo.On("Get", mock.Anything, mock.Anything).Return(expectedResult, nil)
	serviceT := NewService(repo)
	ctx := context.Background()
	result, _ := serviceT.Get(ctx, 2)
	assert.Equal(t, expectedResult, result)
}

func TestUpdateExistentWarehouse(t *testing.T) {
	warehouseRequest := domain.Warehouse{
		ID:                 2,
		Address:            "Monroe 860",
		Telephone:          "47470000",
		WarehouseCode:      "CAD",
		MinimumCapacity:    50,
		MinimumTemperature: 10,
	}
	repo := new(repoM)
	repo.On("Update", mock.Anything, mock.Anything).Return(nil)
	serviceT := NewService(repo)
	ctx := context.Background()
	err := serviceT.Update(ctx, warehouseRequest)
	assert.Nil(t, err)
}

func TestUpdateNonExistentWarehouse(t *testing.T) {
	repo := new(repoM)
	repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("warehouse id not found"))
	serviceT := NewService(repo)
	ctx := context.Background()
	err := serviceT.Update(ctx, domain.Warehouse{})
	assert.NotNil(t, err)
}

func TestDeleteNonExistentWarehouse(t *testing.T) {
	repo := new(repoM)
	repo.On("Delete", mock.Anything, mock.Anything).Return(errors.New("warehouse id not found"))
	serviceT := NewService(repo)
	ctx := context.Background()
	err := serviceT.Delete(ctx, 1)
	assert.NotNil(t, err)
}

func TestDeleteOK(t *testing.T) {
	repo := new(repoM)
	repo.On("Delete", mock.Anything, mock.Anything).Return(nil)
	serviceT := NewService(repo)
	ctx := context.Background()
	err := serviceT.Delete(ctx, 1)
	assert.Nil(t, err)
}

func TestCreateWarehouseError(t *testing.T) {
	warehouseRequest := domain.Warehouse{
		Address:            "Monroe 860",
		Telephone:          "47470000",
		WarehouseCode:      "CAD",
		MinimumCapacity:    50,
		MinimumTemperature: 10,
	}
	repo := new(repoM)
	repo.On("Save", mock.Anything, mock.Anything).Return(1, errors.New("Error al crear el warehouse"))
	serviceT := NewService(repo)
	ctx := context.Background()
	_, err := serviceT.Save(ctx, warehouseRequest)
	assert.NotNil(t, err)
}
