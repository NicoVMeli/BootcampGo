package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/product"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

// Mock interface
func (m *ServiceMock) GetAll(ctx context.Context) ([]domain.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *ServiceMock) Get(ctx context.Context, id int) (domain.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *ServiceMock) Save(ctx context.Context, product domain.Product) (domain.Product, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *ServiceMock) Update(ctx context.Context, id int, product domain.Product) (domain.Product, error) {
	args := m.Called(ctx, id, product)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *ServiceMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ServiceMock) GetRecordReportsByProductId(ctx context.Context, productId int) (string, int, error) {
	args := m.Called(ctx, productId)
	return "", args.Int(0), args.Error(1)
}

// GIN test server
func StartServer(productHandler *Product) *gin.Engine {
	router := gin.Default()
	productRouter := router.Group("api/v1/products")
	{
		productRouter.GET("/", productHandler.GetAll())
		productRouter.GET("/:id", productHandler.Get())
		productRouter.POST("/", productHandler.Create())
		productRouter.PATCH("/:id", productHandler.Update())
		productRouter.DELETE("/:id", productHandler.Delete())
	}
	return router
}

// Create a request for a test
func createTestRequest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

// TESTS

func TestCreateOk(t *testing.T) {
	// ARRANGE
	type response struct {
		Data domain.Product `json:"data"`
	}

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

	createdProduct := domain.Product{
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

	productService := new(ServiceMock)
	productService.On("Save", mock.Anything, productToCreate).Return(createdProduct, nil)

	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodPost, "/api/v1/products/", `{"description": "Yogurt","expiration_rate": 1,"freezing_rate": 2,"height": 6.4,"length": 4.5,"netweight": 3.4,"product_code": "PROD08","recommended_freezing_temperature": 1.3,"width": 1.2,"product_type_id": 2,"seller_id": 4}`)
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(response)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 201, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, createdProduct, testResponse.Data)
}

func TestCreateFailProductCodeRequired(t *testing.T) {
	// ARRANGE
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	expectedCode := "unprocessable_entity"
	expectedMessage := "error: product code is required"

	productService := new(ServiceMock)
	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodPost, "/api/v1/products/", `{"description": "Yogurt","expiration_rate": 1,"freezing_rate": 2,"height": 6.4,"length": 4.5,"netweight": 3.4,"recommended_freezing_temperature": 1.3,"width": 1.2,"product_type_id": 2,"seller_id": 4}`)
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(errorResponse)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 422, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, testResponse.Code)
	assert.Equal(t, expectedMessage, testResponse.Message)
}

func TestCreateFailDescriptionRequired(t *testing.T) {
	// ARRANGE
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	expectedCode := "unprocessable_entity"
	expectedMessage := "error: product description is required"

	productService := new(ServiceMock)
	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodPost, "/api/v1/products/", `{"expiration_rate": 1,"freezing_rate": 2,"height": 6.4,"length": 4.5,"netweight": 3.4,"product_code": "PROD08","recommended_freezing_temperature": 1.3,"width": 1.2,"product_type_id": 2,"seller_id": 4}`)
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(errorResponse)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 422, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, testResponse.Code)
	assert.Equal(t, expectedMessage, testResponse.Message)
}

func TestCreateFailExpirationRateRequired(t *testing.T) {
	// ARRANGE
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	expectedCode := "unprocessable_entity"
	expectedMessage := "error: product expiration rate is required"

	productService := new(ServiceMock)
	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodPost, "/api/v1/products/", `{"description": "Yogurt","freezing_rate": 2,"height": 6.4,"length": 4.5,"netweight": 3.4,"product_code": "PROD08","recommended_freezing_temperature": 1.3,"width": 1.2,"product_type_id": 2,"seller_id": 4}`)
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(errorResponse)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 422, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, testResponse.Code)
	assert.Equal(t, expectedMessage, testResponse.Message)
}

func TestCreateFailFreezingRateRequired(t *testing.T) {
	// ARRANGE
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	expectedCode := "unprocessable_entity"
	expectedMessage := "error: product freezing rate is required"

	productService := new(ServiceMock)
	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodPost, "/api/v1/products/", `{"description": "Yogurt","expiration_rate": 1,"height": 6.4,"length": 4.5,"netweight": 3.4,"product_code": "PROD08","recommended_freezing_temperature": 1.3,"width": 1.2,"product_type_id": 2,"seller_id": 4}`)
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(errorResponse)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 422, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, testResponse.Code)
	assert.Equal(t, expectedMessage, testResponse.Message)
}

func TestCreateFailHeightRequired(t *testing.T) {
	// ARRANGE
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	expectedCode := "unprocessable_entity"
	expectedMessage := "error: product height is required"

	productService := new(ServiceMock)
	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodPost, "/api/v1/products/", `{"description": "Yogurt","expiration_rate": 1,"freezing_rate": 2,"length": 4.5,"netweight": 3.4,"product_code": "PROD08","recommended_freezing_temperature": 1.3,"width": 1.2,"product_type_id": 2,"seller_id": 4}`)
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(errorResponse)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 422, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, testResponse.Code)
	assert.Equal(t, expectedMessage, testResponse.Message)
}

func TestCreateFailLengthRequired(t *testing.T) {
	// ARRANGE
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	expectedCode := "unprocessable_entity"
	expectedMessage := "error: product length is required"

	productService := new(ServiceMock)
	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodPost, "/api/v1/products/", `{"description": "Yogurt","expiration_rate": 1,"freezing_rate": 2,"height": 4.6,"netweight": 3.4,"product_code": "PROD08","recommended_freezing_temperature": 1.3,"width": 1.2,"product_type_id": 2,"seller_id": 4}`)
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(errorResponse)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 422, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, testResponse.Code)
	assert.Equal(t, expectedMessage, testResponse.Message)
}

func TestCreateFailNetweightRequired(t *testing.T) {
	// ARRANGE
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	expectedCode := "unprocessable_entity"
	expectedMessage := "error: product netweight is required"

	productService := new(ServiceMock)
	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodPost, "/api/v1/products/", `{"description": "Yogurt","expiration_rate": 1,"freezing_rate": 2,"height": 4.6,"length": 4.5,"product_code": "PROD08","recommended_freezing_temperature": 1.3,"width": 1.2,"product_type_id": 2,"seller_id": 4}`)
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(errorResponse)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 422, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, testResponse.Code)
	assert.Equal(t, expectedMessage, testResponse.Message)
}

func TestCreateFailRecommendedFreezingTemperatureRequired(t *testing.T) {
	// ARRANGE
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	expectedCode := "unprocessable_entity"
	expectedMessage := "error: product recommended freezing temperature is required"

	productService := new(ServiceMock)
	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodPost, "/api/v1/products/", `{"description": "Yogurt","expiration_rate": 1,"freezing_rate": 2,"height": 4.6,"length": 4.5,"netweight": 3.4,"product_code": "PROD08","width": 1.2,"product_type_id": 2,"seller_id": 4}`)
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(errorResponse)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 422, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, testResponse.Code)
	assert.Equal(t, expectedMessage, testResponse.Message)
}

func TestCreateConflict(t *testing.T) {
	// ARRANGE
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	expectedCode := "conflict"
	expectedMessage := "error: product already exists"

	productService := new(ServiceMock)
	productService.On("Save", mock.Anything, mock.AnythingOfType("domain.Product")).Return(domain.Product{}, product.ErrAlreadyExists)

	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodPost, "/api/v1/products/", `{"description": "Yogurt","expiration_rate": 1,"freezing_rate": 2,"height": 6.4,"length": 4.5,"netweight": 3.4,"product_code": "PROD08","recommended_freezing_temperature": 1.3,"width": 1.2,"product_type_id": 2,"seller_id": 4}`)
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(errorResponse)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 409, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, testResponse.Code)
	assert.Equal(t, expectedMessage, testResponse.Message)
}

func TestFindAll(t *testing.T) {
	// ARRANGE
	type response struct {
		Data []domain.Product `json:"data"`
	}

	var mockedServiceResult []domain.Product
	mockedServiceResult = append(mockedServiceResult, domain.Product{
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

	productService := new(ServiceMock)
	productService.On("GetAll", mock.Anything).Return(mockedServiceResult, nil)

	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodGet, "/api/v1/products/", "")
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(response)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 200, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(testResponse.Data))
	assert.Equal(t, mockedServiceResult, testResponse.Data)
}

func TestFindAllInternalDBError(t *testing.T) {
	// ARRANGE
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	expectedCode := "internal_server_error"
	expectedMessage := "error: internal DB error"

	productService := new(ServiceMock)
	productService.On("GetAll", mock.Anything).Return([]domain.Product{}, errors.New("internal DB error"))

	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodGet, "/api/v1/products/", "")
	productRouter.ServeHTTP(responseRecorder, request)
	assert.Equal(t, 500, responseRecorder.Code)
	//fmt.Println(responseRecorder.Body)
	testResponse := new(errorResponse)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT

	assert.Nil(t, err)
	assert.Equal(t, expectedCode, testResponse.Code)
	assert.Equal(t, expectedMessage, testResponse.Message)
}

func TestFindByIdNonExistent(t *testing.T) {
	// ARRANGE
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	expectedCode := "not_found"
	expectedMessage := "error: product not found"

	productService := new(ServiceMock)
	productService.On("Get", mock.Anything, mock.AnythingOfType("int")).Return(domain.Product{}, product.ErrNotFound)

	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodGet, "/api/v1/products/1", ``)
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(errorResponse)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 404, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, testResponse.Code)
	assert.Equal(t, expectedMessage, testResponse.Message)
}

func TestFindByInvalidId(t *testing.T) {
	// ARRANGE
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	expectedCode := "bad_request"
	expectedMessage := "error: provided id 'ASD' is not an integer"

	productService := new(ServiceMock)
	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodGet, "/api/v1/products/ASD", ``)
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(errorResponse)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 400, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, testResponse.Code)
	assert.Equal(t, expectedMessage, testResponse.Message)
}

func TestFindByIdExistent(t *testing.T) {
	// ARRANGE
	type response struct {
		Data domain.Product `json:"data"`
	}

	mockedServiceResult := domain.Product{
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

	productService := new(ServiceMock)
	productService.On("Get", mock.Anything, 1).Return(mockedServiceResult, nil)

	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodGet, "/api/v1/products/1", "")
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(response)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 200, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, mockedServiceResult, testResponse.Data)
}

func TestUpdateOk(t *testing.T) {
	// ARRANGE
	type response struct {
		Data domain.Product `json:"data"`
	}

	mockedServiceResult := domain.Product{
		ID:             1,
		Description:    "New Yogurt",
		ExpirationRate: 3,
		FreezingRate:   4,
		Height:         6.4,
		Length:         4.5,
		Netweight:      3.4,
		ProductCode:    "NEWPROD08",
		RecomFreezTemp: 1.3,
		Width:          1.2,
		ProductTypeID:  2,
		SellerID:       4,
	}

	productService := new(ServiceMock)
	productService.On("Update", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("domain.Product")).Return(mockedServiceResult, nil)

	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodPatch, "/api/v1/products/1", `{"description": "New Yogurt","expiration_rate": 3,"freezing_rate": 4,"product_code": "NEWPROD08"}`)
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(response)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 200, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, mockedServiceResult, testResponse.Data)
}

func TestUpdateNonExistent(t *testing.T) {
	// ARRANGE
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	expectedCode := "not_found"
	expectedMessage := "error: product not found"

	productService := new(ServiceMock)
	productService.On("Update", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("domain.Product")).Return(domain.Product{}, product.ErrNotFound)

	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodPatch, "/api/v1/products/1", `{"description": "New Yogurt","expiration_rate": 3,"freezing_rate": 4,"product_code": "NEWPROD08"}`)
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(errorResponse)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 404, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, testResponse.Code)
	assert.Equal(t, expectedMessage, testResponse.Message)
}

func TestDeleteNonExistent(t *testing.T) {
	// ARRANGE
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	expectedCode := "not_found"
	expectedMessage := "error: product not found"

	productService := new(ServiceMock)
	productService.On("Delete", mock.Anything, mock.AnythingOfType("int")).Return(product.ErrNotFound)

	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodDelete, "/api/v1/products/1", "")
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(errorResponse)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 404, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, testResponse.Code)
	assert.Equal(t, expectedMessage, testResponse.Message)
}

func TestDeleteOk(t *testing.T) {
	// ARRANGE
	type response struct {
		Data []domain.Product `json:"data"`
	}

	productService := new(ServiceMock)
	productService.On("Delete", mock.Anything, mock.AnythingOfType("int")).Return(nil)

	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodDelete, "/api/v1/products/1", "")
	productRouter.ServeHTTP(responseRecorder, request)

	// ASSERT
	assert.Equal(t, 204, responseRecorder.Code)
	assert.Equal(t, 0, responseRecorder.Body.Len())
}

func TestDeleteByInvalidId(t *testing.T) {
	// ARRANGE
	type errorResponse struct {
		Status  int    `json:"-"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	expectedCode := "bad_request"
	expectedMessage := "error: provided id 'ASD' is not an integer"

	productService := new(ServiceMock)
	productHandler := NewProduct(productService)
	productRouter := StartServer(productHandler)

	// ACT
	request, responseRecorder := createTestRequest(http.MethodDelete, "/api/v1/products/ASD", ``)
	productRouter.ServeHTTP(responseRecorder, request)

	testResponse := new(errorResponse)
	err := json.Unmarshal(responseRecorder.Body.Bytes(), testResponse)

	// ASSERT
	assert.Equal(t, 400, responseRecorder.Code)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, testResponse.Code)
	assert.Equal(t, expectedMessage, testResponse.Message)
}
