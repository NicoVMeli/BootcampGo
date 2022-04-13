package employee

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type dbMock struct {
	mock.Mock
}

func (m *dbMock) GetAll(ctx context.Context) ([]domain.Employee, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Employee), args.Error(1)
}

func (m *dbMock) Get(ctx context.Context, id int) (domain.Employee, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Employee), args.Error(1)
}

func (m *dbMock) Exists(ctx context.Context, cardNumberID string) bool {
	args := m.Called(ctx, cardNumberID)
	return args.Bool(0)
}

func (m *dbMock) Save(ctx context.Context, e domain.Employee) (int, error) {
	args := m.Called(ctx, e)
	return args.Int(0), args.Error(1)
}

func (m *dbMock) Update(ctx context.Context, e domain.Employee) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *dbMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *dbMock) GetAllReportInboundOrders() ([]domain.ReportInboundOrders, error) {
	args := m.Called()
	return args.Get(0).([]domain.ReportInboundOrders), args.Error(1)
}

func (m *dbMock) GetReportByEmployeeIdInboundOrders(employeeId int) (domain.ReportInboundOrders, error) {
	args := m.Called(employeeId)
	return args.Get(0).(domain.ReportInboundOrders), args.Error(1)
}

func TestCreateOk(t *testing.T) {
	expectedResult := domain.Employee{
		CardNumberID: "123456",
		FirstName:    "laynerker",
		LastName:     "guerrero",
		WarehouseID:  5,
	}
	repo := new(dbMock)
	repo.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	repo.On("Exists", mock.Anything, mock.Anything).Return(false)
	serviceT := NewService(repo)
	ctx := context.Background()
	resul, _ := serviceT.Save(ctx, expectedResult)
	assert.True(t, resul.ID == 1)
}

func TestCreateConflict(t *testing.T) {
	expectedError := errors.New("El usuario ya existe")
	expectedResult := domain.Employee{}
	repo := new(dbMock)
	repo.On("Exists", mock.Anything, mock.Anything).Return(true)
	serviceT := NewService(repo)
	ctx := context.Background()
	resul, err := serviceT.Save(ctx, expectedResult)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, expectedResult, resul)
}

func TestFindAll(t *testing.T) {
	expectedResult := []domain.Employee{}
	expectedResult = append(expectedResult, domain.Employee{
		ID:           1,
		CardNumberID: "123456",
		FirstName:    "laynerker",
		LastName:     "guerrero",
		WarehouseID:  5,
	})
	repo := new(dbMock)
	repo.On("GetAll", mock.Anything).Return(expectedResult, nil)
	serviceT := NewService(repo)
	ctx := context.Background()
	resul, _ := serviceT.GetAll(ctx)
	assert.True(t, len(resul) > 0)
}

func TestFindByIdExistent(t *testing.T) {
	expectedResult := domain.Employee{
		ID:           1,
		CardNumberID: "123456",
		FirstName:    "laynerker",
		LastName:     "guerrero",
		WarehouseID:  5,
	}
	repo := new(dbMock)
	repo.On("Get", mock.Anything, mock.Anything).Return(expectedResult, nil)
	serviceT := NewService(repo)
	ctx := context.Background()
	resul, _ := serviceT.Get(ctx, 1)
	assert.Equal(t, expectedResult, resul)
}

func TestFindByIdNonExistent(t *testing.T) {
	expectedError := fmt.Errorf("No existe el empleado con el id %d", 1)
	expectedResult := domain.Employee{}
	repo := new(dbMock)
	repo.On("Get", mock.Anything, mock.Anything).Return(expectedResult, errors.New("sql: no rows in result set"))
	serviceT := NewService(repo)
	ctx := context.Background()
	result, err := serviceT.Get(ctx, 1)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, expectedResult, result)
}

func TestUpdateExistent(t *testing.T) {
	expectedResultGet := domain.Employee{
		ID:           1,
		CardNumberID: "123456",
		FirstName:    "laynerker",
		LastName:     "guerrero",
		WarehouseID:  5,
	}
	expectedResult := domain.Employee{
		ID:           1,
		CardNumberID: "12345",
		FirstName:    "laynerke",
		LastName:     "guerrer",
		WarehouseID:  6,
	}
	repo := new(dbMock)
	repo.On("Get", mock.Anything, mock.Anything).Return(expectedResultGet, nil)
	repo.On("Update", mock.Anything, mock.Anything).Return(nil)
	serviceT := NewService(repo)
	ctx := context.Background()
	result, _ := serviceT.Update(ctx, 1, expectedResult)
	assert.Equal(t, expectedResult, result)
}

func TestUpdateNonExistent(t *testing.T) {
	expectedResult := domain.Employee{}
	expectedError := fmt.Errorf("No existe el empleado con el id %d", 1)
	repo := new(dbMock)
	repo.On("Get", mock.Anything, mock.Anything).Return(expectedResult, errors.New("El usuario no existe"))
	serviceT := NewService(repo)
	ctx := context.Background()
	result, err := serviceT.Update(ctx, 1, expectedResult)
	assert.Equal(t, expectedResult, result)
	assert.Equal(t, expectedError, err)
}

func TestDeleteNonExistent(t *testing.T) {
	expectedError := fmt.Errorf("No existe el empleado con el id %d", 1)
	expectedResult := domain.Employee{}
	repo := new(dbMock)
	repo.On("Get", mock.Anything, mock.Anything).Return(expectedResult, errors.New("El usuario no existe"))
	serviceT := NewService(repo)
	ctx := context.Background()
	err := serviceT.Delete(ctx, 1)
	assert.Equal(t, expectedError, err)
}

func TestDeleteOK(t *testing.T) {
	expectedResultGet := domain.Employee{
		ID:           1,
		CardNumberID: "123456",
		FirstName:    "laynerker",
		LastName:     "guerrero",
		WarehouseID:  5,
	}
	repo := new(dbMock)
	repo.On("Get", mock.Anything, mock.Anything).Return(expectedResultGet, nil)
	repo.On("Delete", mock.Anything, mock.Anything).Return(nil)
	serviceT := NewService(repo)
	ctx := context.Background()
	err := serviceT.Delete(ctx, 1)
	assert.Nil(t, err)
}

func TestGetAllReportInboundOrdersOK(t *testing.T) {
	expectedResult := []domain.ReportInboundOrders{}
	expectedResult = append(expectedResult, domain.ReportInboundOrders{
		ID:                 1,
		CardNumberID:       "123456",
		FirstName:          "laynerker",
		LastName:           "guerrero",
		WarehouseID:        5,
		InboundOrdersCount: 3,
	})
	repo := new(dbMock)
	repo.On("GetAllReportInboundOrders").Return(expectedResult, nil)
	serviceT := NewService(repo)
	ctx := context.Background()
	resul, _ := serviceT.GetAllReportInboundOrders(ctx)
	assert.True(t, len(resul) > 0)
}

func TestGetReportByEmployeeIdInboundOrdersOK(t *testing.T) {
	expectedResult := domain.ReportInboundOrders{
		ID:                 1,
		CardNumberID:       "123456",
		FirstName:          "laynerker",
		LastName:           "guerrero",
		WarehouseID:        5,
		InboundOrdersCount: 2,
	}
	repo := new(dbMock)
	repo.On("GetReportByEmployeeIdInboundOrders", mock.Anything).Return(expectedResult, nil)
	serviceT := NewService(repo)
	ctx := context.Background()
	resul, _ := serviceT.GetReportByEmployeeIdInboundOrders(ctx, 1)
	assert.Equal(t, expectedResult, resul)
}

func TestGetReportByEmployeeIdInboundOrdersErrorQuery(t *testing.T) {
	expectedResult := domain.ReportInboundOrders{}
	expectedResultError := errors.New("The employee doesn't exist")
	repo := new(dbMock)
	repo.On("GetReportByEmployeeIdInboundOrders", mock.Anything).Return(expectedResult, errors.New("Sql Error"))
	serviceT := NewService(repo)
	ctx := context.Background()
	resul, err := serviceT.GetReportByEmployeeIdInboundOrders(ctx, 1)
	assert.Equal(t, expectedResult, resul)
	assert.Equal(t, expectedResultError, err)
}
