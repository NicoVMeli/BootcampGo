package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/employee"
	"github.com/gin-gonic/gin"
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

func createServiceE(p *Employee) *gin.Engine {
	r := gin.Default()
	pr := r.Group("api/v1/employees")
	{
		pr.GET("/", p.GetAll())
		pr.GET("/:id", p.Get())
		pr.GET("/reportinboundorders", p.GetReportIO())
		pr.POST("/", p.Create())
		pr.PATCH("/:id", p.Update())
		pr.DELETE("/:id", p.Delete())
	}
	return r
}

func createRequestEmployeeTest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

func TestFindAllEmployee(t *testing.T) {
	type response struct {
		Data []domain.Employee
	}
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
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodGet, "/api/v1/employees/", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	var objRes response
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, objRes.Data)
}

func TestGetAllNoRecordsExist(t *testing.T) {
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	expectedError := "No Existen Registros"
	expectedResult := []domain.Employee{}
	repo := new(dbMock)
	repo.On("GetAll", mock.Anything).Return(expectedResult, nil)
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodGet, "/api/v1/employees/", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	var objRes errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedError, objRes.Message)
}

func TestFindByIdNonExistentEmployee(t *testing.T) {
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	expectedError := "No existe el empleado con el id 456"
	expectedResult := domain.Employee{}
	repo := new(dbMock)
	repo.On("Get", mock.Anything, mock.Anything).Return(expectedResult, errors.New("sql: no rows in result set"))
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodGet, "/api/v1/employees/456", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	var objRes errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedError, objRes.Message)
}

func TestFindByIdExistentEmployee(t *testing.T) {
	type response struct {
		Data domain.Employee
	}
	expectedResult := domain.Employee{
		ID:           1,
		CardNumberID: "123456",
		FirstName:    "laynerker",
		LastName:     "guerrero",
		WarehouseID:  5,
	}
	repo := new(dbMock)
	repo.On("Get", mock.Anything, mock.Anything).Return(expectedResult, nil)
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodGet, "/api/v1/employees/1", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	var objRes response
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, objRes.Data)
}

func TestGetInvalidId(t *testing.T) {
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	expectedError := "Error: invalid ID"
	repo := new(dbMock)
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodGet, "/api/v1/employees/ASD", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var objRes errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedError, objRes.Message)
}

func TestCreateOkEmployee(t *testing.T) {
	type response struct {
		Data domain.Employee
	}
	expectedResult := domain.Employee{
		ID:           1,
		CardNumberID: "123456",
		FirstName:    "laynerker",
		LastName:     "guerrero",
		WarehouseID:  5,
	}
	repo := new(dbMock)
	repo.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	repo.On("Exists", mock.Anything, mock.Anything).Return(false)
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodPost, "/api/v1/employees/", `{"card_number_id": "123456","first_name": "laynerker","last_name": "guerrero","warehouse_id": 5}`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var objRes response
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, objRes.Data)
}

func TestCreateFail(t *testing.T) {
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	expectedError := "Todos los campos son requeridos"
	repo := new(dbMock)
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodPost, "/api/v1/employees/", `{"card_number_id": "123456"}`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	var objRes errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedError, objRes.Message)
}

func TestCreateConflictEmployee(t *testing.T) {
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	expectedError := "El usuario ya existe"
	repo := new(dbMock)
	repo.On("Exists", mock.Anything, mock.Anything).Return(true)
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodPost, "/api/v1/employees/", `{"card_number_id": "123456","first_name": "laynerker","last_name": "guerrero","warehouse_id": 5}`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	var objRes errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedError, objRes.Message)
}

func TestUpdateOK(t *testing.T) {
	type response struct {
		Data domain.Employee
	}
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
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodPatch, "/api/v1/employees/1", `{"card_number_id": "12345","first_name": "laynerke","last_name": "guerrer","warehouse_id": 6}`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	var objRes response
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, objRes.Data)
}

func TestUpdateNonExistentEmployee(t *testing.T) {
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	expectedResult := domain.Employee{}
	expectedError := "Error in code for: No existe el empleado con el id 1"
	repo := new(dbMock)
	repo.On("Get", mock.Anything, mock.Anything).Return(expectedResult, errors.New("El usuario no existe"))
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodPatch, "/api/v1/employees/1", `{"card_number_id": "123456","first_name": "laynerker","last_name": "guerrero","warehouse_id": 5}`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	var objRes errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedError, objRes.Message)
}

func TestDeleteNonExistentEmployee(t *testing.T) {
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	expectedResult := domain.Employee{}
	expectedError := "Error in code for: No existe el empleado con el id 1"
	repo := new(dbMock)
	repo.On("Get", mock.Anything, mock.Anything).Return(expectedResult, errors.New("El usuario no existe"))
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodDelete, "/api/v1/employees/1", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	var objRes errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedError, objRes.Message)
}

func TestDeleteOK(t *testing.T) {
	type response struct {
		Data domain.Employee
	}
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
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodDelete, "/api/v1/employees/1", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteInvalidId(t *testing.T) {
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	expectedError := "Error: invalid ID"
	repo := new(dbMock)
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodDelete, "/api/v1/employees/ASD", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var objRes errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedError, objRes.Message)
}

func TestReportGetAllInboudOrders(t *testing.T) {
	type response struct {
		Data []domain.ReportInboundOrders
	}
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
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodGet, "/api/v1/employees/reportinboundorders", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	var objRes response
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, objRes.Data)
}

func TestReportGetAllInboudOrdersNoResult(t *testing.T) {
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	expectedResult := []domain.ReportInboundOrders{}
	expectedError := "No existen resultados"
	repo := new(dbMock)
	repo.On("GetAllReportInboundOrders").Return(expectedResult, errors.New(expectedError))
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodGet, "/api/v1/employees/reportinboundorders", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	var objRes errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("Error in code for: %s", expectedError), objRes.Message)
}

func TestInboudOrdersReportGetByEmployee(t *testing.T) {
	type response struct {
		Data domain.ReportInboundOrders
	}
	expectedResult := domain.ReportInboundOrders{
		ID:                 1,
		CardNumberID:       "123456",
		FirstName:          "laynerker",
		LastName:           "guerrero",
		WarehouseID:        5,
		InboundOrdersCount: 3,
	}
	repo := new(dbMock)
	repo.On("GetReportByEmployeeIdInboundOrders", mock.Anything).Return(expectedResult, nil)
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodGet, "/api/v1/employees/reportinboundorders?id=1", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	var objRes response
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, objRes.Data)
}

func TestGetReportInboudOrdersInvalidId(t *testing.T) {
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	expectedError := "Error: invalid ID"
	repo := new(dbMock)
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodGet, "/api/v1/employees/reportinboundorders?id=ASD", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var objRes errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedError, objRes.Message)
}

func TestGetReportInboudOrdersNoResult(t *testing.T) {
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	id := 1
	expectedError := fmt.Sprintf("No existen resultados para el empleado con el id %d", id)
	repo := new(dbMock)
	repo.On("GetReportByEmployeeIdInboundOrders", mock.Anything).Return(domain.ReportInboundOrders{}, errors.New("The employee doesn't exist"))
	service := employee.NewService(repo)
	p := NewEmployee(service)
	r := createServiceE(p)
	req, rr := createRequestTest(http.MethodGet, fmt.Sprintf("/api/v1/employees/reportinboundorders?id=%d", id), "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	var objRes errorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &objRes)
	assert.Nil(t, err)
	assert.Equal(t, expectedError, objRes.Message)
}
