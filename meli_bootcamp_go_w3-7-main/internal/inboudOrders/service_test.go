package inboudOrders

import (
	"context"
	"testing"
	"time"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type dbIOMock struct {
	mock.Mock
}

func (m *dbIOMock) Save(ctx context.Context, i domain.InboudOrders) (int, error) {
	args := m.Called(ctx, i)
	return args.Int(0), args.Error(1)
}

func (m *dbIOMock) ExistsEmployee(ctx context.Context, employeeId int) bool {
	args := m.Called(ctx, employeeId)
	return args.Bool(0)
}

func TestCreateOkInboudOrders(t *testing.T) {
	type response struct {
		Data domain.InboudOrders
	}
	orderDate, _ := time.Parse("2006-01-02", "2022-01-06")
	expectedResult := domain.InboudOrders{
		ID:             1,
		OrderDate:      orderDate,
		OrderNumber:    "12312",
		EmployeeId:     1,
		ProductBatchId: 2,
		WarehouseID:    3,
	}
	repo := new(dbIOMock)
	repo.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	repo.On("ExistsEmployee", mock.Anything, mock.Anything).Return(true)
	serviceIO := NewService(repo)
	ctx := context.Background()
	resul, _ := serviceIO.Save(ctx, expectedResult)
	assert.True(t, resul.ID == 1)
}

func TestCreateInboudOrdersNoExistEmployee(t *testing.T) {
	type response struct {
		Data domain.InboudOrders
	}
	orderDate, _ := time.Parse("2006-01-02", "2022-01-06")
	expectedResult := domain.InboudOrders{
		ID:             1,
		OrderDate:      orderDate,
		OrderNumber:    "12312",
		EmployeeId:     1,
		ProductBatchId: 2,
		WarehouseID:    3,
	}
	repo := new(dbIOMock)
	//repo.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	repo.On("ExistsEmployee", mock.Anything, mock.Anything).Return(false)
	serviceIO := NewService(repo)
	ctx := context.Background()
	_, err := serviceIO.Save(ctx, expectedResult)
	assert.Error(t, err)
}
