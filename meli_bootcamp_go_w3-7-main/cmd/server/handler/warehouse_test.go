package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/warehouse"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ServiceM struct {
	mock.Mock
}

// Mock interface
func (m *ServiceM) Save(ctx context.Context, w domain.Warehouse) (int, error) {
	args := m.Called(ctx, w)
	return args.Int(0), args.Error(1)
}

func (m *ServiceM) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Warehouse), args.Error(1)
}

func (m *ServiceM) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Warehouse), args.Error(1)
}

func (m *ServiceM) Update(ctx context.Context, w domain.Warehouse) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *ServiceM) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ServiceM) Exists(ctx context.Context, warehouseCode string) bool {
	args := m.Called(ctx, warehouseCode)
	return args.Bool(0)
}

func createServer(w *Warehouse) *gin.Engine {
	r := gin.Default()
	wr := r.Group("/api/v1/warehouses")
	{
		wr.GET("/", w.GetAll())
		wr.GET("/:id", w.Get())
		wr.POST("/", w.Create())
		wr.PATCH("/:id", w.Update())
		wr.DELETE("/:id", w.Delete())
	}
	return r
}

func createRequestTest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

// TESTS
func TestCreateWarehouseOK(t *testing.T) {
	expectedResponse := `{"data":{"id":1,"address":"Hernan Cortez 155","telephone":"4995687","warehouse_code":"HCD","minimum_capacity":30,"minimum_temperature":2}}`
	s := new(ServiceM)
	s.On("Exists", mock.Anything, mock.Anything).Return(false)
	s.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	service := warehouse.NewService(s)
	w := NewWarehouse(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodPost, `/api/v1/warehouses/`,
		`{
		"address": "Hernan Cortez 155",
		"telephone": "4995687",
		"warehouse_code": "HCD",
		"minimum_capacity": 30,
		"minimum_temperature": 2
		}`)
	r.ServeHTTP(res, req)
	actualResponse := res.Body.String()
	assert.Equal(t, http.StatusCreated, res.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestCreateWarehouseFail(t *testing.T) {
	s := new(ServiceM)
	s.On("Exists", mock.Anything, mock.Anything).Return(false)
	s.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	service := warehouse.NewService(s)
	w := NewWarehouse(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodPost, `/api/v1/warehouses/`,
		`{
		"address": "Hernan Cortez 155",
		"telephone": "4995687",
		"warehouse_code": "HCD",
		"minimum_capacity": 30
		}`)
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusUnprocessableEntity, res.Code)

}

func TestCreateWarehouseConflict(t *testing.T) {
	s := new(ServiceM)
	s.On("Exists", mock.Anything, mock.Anything).Return(true)
	s.On("Save", mock.Anything, mock.AnythingOfType("domain.Warehouse")).Return(1, errors.New("Ya existe un Warehouse asociado al código que intenta registrar. El campo 'warehouse_code' debe ser único"))
	service := warehouse.NewService(s)
	w := NewWarehouse(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodPost, "/api/v1/warehouses/",
		`{
        "address": "Hernan Cortez 155",
        "telephone": "4995687",
        "warehouse_code": "HCD",
        "minimum_capacity": 30,
		"minimum_temperature": 2
    }`)
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusConflict, res.Code)
}

func TestFindAllWarehouses(t *testing.T) {
	type response struct {
		Data []domain.Warehouse
	}
	mockResponse := []domain.Warehouse{}
	mockResponse = append(mockResponse, domain.Warehouse{
		ID:                 1,
		Address:            "Monroe 860",
		Telephone:          "47470000",
		WarehouseCode:      "CAD",
		MinimumCapacity:    50,
		MinimumTemperature: 10,
	})
	mockResponse = append(mockResponse, domain.Warehouse{
		ID:                 2,
		Address:            "Casa 20",
		Telephone:          "7223828",
		WarehouseCode:      "SHA",
		MinimumCapacity:    10,
		MinimumTemperature: 7,
	})
	s := new(ServiceM)
	s.On("GetAll", mock.Anything, mock.Anything).Return(mockResponse, nil)
	service := warehouse.NewService(s)
	w := NewWarehouse(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodGet, "/api/v1/warehouses/", "")
	r.ServeHTTP(res, req)
	type dataResponse struct {
		data struct {
			id                  int
			address             string
			telephone           string
			warehouse_code      string
			minimun_capacity    int
			minimun_temperature int
		}
	}
	assert.Equal(t, http.StatusOK, res.Code)

}

func TestFindWarehouseByIdNonExistent(t *testing.T) {
	expectedResult := domain.Warehouse{}
	expectedError := "sql: no rows in result set"
	s := new(ServiceM)
	s.On("Get", mock.Anything, mock.Anything).Return(expectedResult, errors.New(expectedError))
	service := warehouse.NewService(s)
	w := NewWarehouse(service)
	r := createServer(w)
	req, err := createRequestTest(http.MethodGet, "/api/v1/warehouses/100", "")
	r.ServeHTTP(err, req)
	assert.Equal(t, http.StatusNotFound, err.Code)
}

func TestFindWarehouseByIdExistent(t *testing.T) {
	mockResponse := domain.Warehouse{
		ID:                 1,
		Address:            "Monroe 860",
		Telephone:          "47470000",
		WarehouseCode:      "CAD",
		MinimumCapacity:    50,
		MinimumTemperature: 10,
	}
	expectedResponse := `{"data":{"id":1,"address":"Monroe 860","telephone":"47470000","warehouse_code":"CAD","minimum_capacity":50,"minimum_temperature":10}}`
	s := new(ServiceM)
	s.On("Get", mock.Anything, mock.Anything).Return(mockResponse, nil)
	service := warehouse.NewService(s)
	w := NewWarehouse(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodGet, "/api/v1/warehouses/1", "")
	r.ServeHTTP(res, req)
	actualResponse := res.Body.String()
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUpdateWarehouseOK(t *testing.T) {
	mockResponse := domain.Warehouse{
		ID:                 2,
		Address:            "Monroe 860",
		Telephone:          "47470000",
		WarehouseCode:      "CAD",
		MinimumCapacity:    50,
		MinimumTemperature: 10,
	}
	expectedResponse := `{"data":{"id":2,"address":"Hernan Cortez 155","telephone":"4995687","warehouse_code":"HCD","minimum_capacity":50,"minimum_temperature":10}}`
	s := new(ServiceM)
	s.On("Get", mock.Anything, mock.Anything).Return(mockResponse, nil)
	s.On("Update", mock.Anything, mock.Anything).Return(nil)
	service := warehouse.NewService(s)
	w := NewWarehouse(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodPatch, `/api/v1/warehouses/2`,
		`{
		"address": "Hernan Cortez 155",
		"telephone": "4995687",
		"warehouse_code": "HCD"
		}`)
	r.ServeHTTP(res, req)
	actualResponse := res.Body.String()
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestUpdateNonExistWarehouse(t *testing.T) {

	expectedResult := domain.Warehouse{}
	expectedError := "sql: no rows in result set"
	s := new(ServiceM)
	s.On("Get", mock.Anything, mock.Anything).Return(expectedResult, errors.New(expectedError))
	service := warehouse.NewService(s)
	w := NewWarehouse(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodPatch, `/api/v1/warehouses/2`,
		`{
		"address": "Hernan Cortez 155",
		"telephone": "4995687",
		"warehouse_code": "HCD"
		}`)
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestDeleteNonExistWarehouse(t *testing.T) {
	expectedError := "sql: no rows in result set"
	s := new(ServiceM)
	s.On("Get", mock.Anything, mock.Anything).Return(domain.Warehouse{}, errors.New(expectedError))
	service := warehouse.NewService(s)
	w := NewWarehouse(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodDelete, `/api/v1/warehouses/2`, "")
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestDeleteWarehouseOK(t *testing.T) {
	mockResponse := domain.Warehouse{
		ID:                 2,
		Address:            "Monroe 860",
		Telephone:          "47470000",
		WarehouseCode:      "CAD",
		MinimumCapacity:    50,
		MinimumTemperature: 10,
	}
	s := new(ServiceM)
	s.On("Get", mock.Anything, mock.Anything).Return(mockResponse, nil)
	s.On("Delete", mock.Anything, mock.Anything).Return(nil)
	service := warehouse.NewService(s)
	w := NewWarehouse(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodDelete, `/api/v1/warehouses/1`, "")
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusNoContent, res.Code)
}

func TestUpdateWarehouseFail(t *testing.T) {
	mockResponse := domain.Warehouse{
		ID:                 2,
		Address:            "Monroe 860",
		Telephone:          "47470000",
		WarehouseCode:      "CAD",
		MinimumCapacity:    50,
		MinimumTemperature: 10,
	}
	s := new(ServiceM)
	s.On("Get", mock.Anything, mock.Anything).Return(mockResponse, nil)
	s.On("Update", mock.Anything, mock.Anything).Return(nil)
	service := warehouse.NewService(s)
	w := NewWarehouse(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodPatch, `/api/v1/warehouses/2`,
		`{
		"address": "Hernan Cortez 155",
		"telephone": "4995687",
		"warehouse_code": "HCD",
		"id": "34"
		}`)
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusBadRequest, res.Code)
}
