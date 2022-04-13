package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/inboudOrders"
	"github.com/gin-gonic/gin"
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

func createServiceIO(p *InboudOrders) *gin.Engine {
	r := gin.Default()
	pr := r.Group("api/v1/inboundOrders")
	{
		pr.POST("/", p.Create())
	}
	return r
}

func createRequestInboudOrdersTest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	return req, httptest.NewRecorder()
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
	service := inboudOrders.NewService(repo)
	p := NewInboudOrders(service)
	r := createServiceIO(p)
	req, rr := createRequestInboudOrdersTest(http.MethodPost, "/api/v1/inboundOrders/", `{
		"order_date": "2022-01-06",
		"order_number": "12312",
		"employee_id": 1,
		"product_batch_id": 2,
		"warehouse_id": 3
		}`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var objRes response
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, objRes.Data)
}

func TestCreateFailInboudOrders(t *testing.T) {
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	expectedError := "El empleado no existe"
	repo := new(dbIOMock)
	service := inboudOrders.NewService(repo)
	p := NewInboudOrders(service)
	r := createServiceIO(p)
	repo.On("ExistsEmployee", mock.Anything, mock.Anything).Return(false)
	req, rr := createRequestInboudOrdersTest(http.MethodPost, "/api/v1/inboundOrders/", `{
		"order_date": "2022-01-06",
		"order_number": "12312",
		"employee_id": 1,
		"product_batch_id": 2,
		"warehouse_id": 3
		}`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusConflict, rr.Code)
	var objRes errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedError, objRes.Message)
}

func TestCreateFailFieldInboudOrders(t *testing.T) {
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	expectedError := "Todos los campos son requeridos"
	repo := new(dbIOMock)
	service := inboudOrders.NewService(repo)
	p := NewInboudOrders(service)
	r := createServiceIO(p)
	req, rr := createRequestInboudOrdersTest(http.MethodPost, "/api/v1/inboundOrders/", `{"order_date": "2022-01-06"}`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	var objRes errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedError, objRes.Message)
}
