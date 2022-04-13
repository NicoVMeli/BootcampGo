package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func createSellerServer() *gin.Engine {
	sellerService := NewSellerServiceMock()
	seller := NewSeller(sellerService)
	r := gin.Default()

	sellersGroup := r.Group("/sellers")
	{
		sellersGroup.GET("/", seller.GetAll())
		sellersGroup.GET("/:id", seller.Get())
		sellersGroup.POST("/", seller.Create())
		sellersGroup.PATCH("/:id", seller.Update())
		sellersGroup.DELETE("/:id", seller.Delete())
	}
	return r
}
func createSellerRequestTest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", "1234")
	return req, httptest.NewRecorder()
}

type SellerServiceMock struct {
	Data []byte
}

func NewSellerServiceMock() *SellerServiceMock {
	datosDePrueba := []domain.Seller{
		{
			ID:          1,
			CID:         1,
			CompanyName: "Meli",
			Address:     "Bulnes 10",
			Telephone:   "123456",
		},
		{
			ID:          2,
			CID:         2,
			CompanyName: "Baires Dev",
			Address:     " Av Belgrano 1200",
			Telephone:   "4833269",
		},
		{
			ID:          3,
			CID:         3,
			CompanyName: "Uala",
			Address:     " Av Mate de Luna 2200",
			Telephone:   "3814471789",
		},
	}
	datos, _ := json.Marshal(datosDePrueba)
	return &SellerServiceMock{
		Data: datos,
	}
}

func (s *SellerServiceMock) GetAll(ctx context.Context) ([]domain.Seller, error) {
	var datos []domain.Seller
	err := json.Unmarshal(s.Data, &datos)
	if err != nil {
		return []domain.Seller{}, err
	}
	return datos, nil
}

func (s *SellerServiceMock) Get(ctx context.Context, id int) (domain.Seller, error) {
	var datos []domain.Seller
	err := json.Unmarshal(s.Data, &datos)
	if err != nil {
		return domain.Seller{}, err
	}
	for _, seller := range datos {
		if seller.ID == id {
			return seller, nil
		}
	}
	return domain.Seller{}, fmt.Errorf("Seller not found")
}

func (s *SellerServiceMock) Save(ctx context.Context, se domain.Seller) (domain.Seller, error) {
	var datos []domain.Seller
	err := json.Unmarshal(s.Data, &datos)
	if err != nil {
		return domain.Seller{}, err
	}
	for _, sel := range datos {
		if sel.CID == se.CID {
			return domain.Seller{}, fmt.Errorf("Seller already exists")
		}
	}
	se.ID = datos[len(datos)-1].ID
	datos = append(datos, se)
	return se, nil
}

func (s *SellerServiceMock) Update(ctx context.Context, se domain.Seller) (domain.Seller, error) {
	var datos []domain.Seller
	err := json.Unmarshal(s.Data, &datos)
	if err != nil {
		return domain.Seller{}, err
	}
	for _, sel := range datos {
		if sel.ID == se.ID {
			sel.CID = se.CID
			sel.CompanyName = se.CompanyName
			sel.Address = se.Address
			sel.Telephone = se.Telephone
			return sel, nil
		}
	}
	return domain.Seller{}, fmt.Errorf("Seller does not exist")
}

func (s *SellerServiceMock) Delete(ctx context.Context, id int) error {
	var datos []domain.Seller
	err := json.Unmarshal(s.Data, &datos)
	if err != nil {
		return err
	}
	for _, seller := range datos {
		if seller.ID == id {
			return nil
		}
	}
	return fmt.Errorf("Seller not found")
}

func (s *SellerServiceMock) Exists(ctx context.Context, cid int) bool {
	return false
}

func Test_create_ok(t *testing.T) {
	r := createSellerServer()

	body := `
		{
			"cid": 4,
			"company_name": "Pedidos Ya",
			"address": "Av Aconquija 1100",
			"telephone": "987654"
		}`
	req, rr := createSellerRequestTest(http.MethodPost, "/sellers/", body)

	r.ServeHTTP(rr, req)

	assert.Equal(t, 201, rr.Code, rr.Result())
}

func Test_create_fail(t *testing.T) {
	r := createSellerServer()
	body := `
		{
			
			"company_name": "Pedidos Ya",
			"address": "Av Aconquija 1100",
			"telephone": "987654"
		}`
	req, rr := createSellerRequestTest(http.MethodPost, "/sellers/", body)

	r.ServeHTTP(rr, req)

	assert.Equal(t, 422, rr.Code, rr.Result())
}

func Test_create_conflict(t *testing.T) {
	r := createSellerServer()

	body := `
		{
			"cid": 1,
			"company_name": "Demo",
			"address": "Bulnes 10",
			"telephone": "123456"
		}`
	req, rr := createSellerRequestTest(http.MethodPost, "/sellers/", body)

	r.ServeHTTP(rr, req)

	assert.Equal(t, 422, rr.Code, rr.Result())
}

func Test_find_all(t *testing.T) {
	r := createSellerServer()
	req, rr := createSellerRequestTest(http.MethodGet, "/sellers/", "{}")

	r.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code, rr.Result())
}

func Test_find_by_id_non_existent(t *testing.T) {
	r := createSellerServer()
	req, rr := createSellerRequestTest(http.MethodGet, "/sellers/4", "{}")
	r.ServeHTTP(rr, req)
	assert.Equal(t, 401, rr.Code, rr.Result())
}

func Test_find_by_id_existent(t *testing.T) {
	r := createSellerServer()
	req, rr := createSellerRequestTest(http.MethodGet, "/sellers/2", "{}")
	r.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code, rr.Result())
}

func Test_update_ok(t *testing.T) {
	r := createSellerServer()

	body := `
		{
			"id": 2,
			"cid": 2,
			"company_name": "Baires Dev",
			"address": "Av Aconquija 1200",
			"telephone": "4833269"
		}`
	req, rr := createSellerRequestTest(http.MethodPatch, "/sellers/2", body)

	r.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}

func Test_update_non_existent(t *testing.T) {
	r := createSellerServer()
	body := `
		{
			"id": 2,
			"cid": 2,
			"company_name": "Baires Dev",
			"address": "Av Aconquija 1200",
			"telephone": "4833269"
		}`
	req, rr := createSellerRequestTest(http.MethodPatch, "/sellers/5", body)

	r.ServeHTTP(rr, req)

	assert.Equal(t, 400, rr.Code)
}

func Test_delete_non_existent(t *testing.T) {
	r := createSellerServer()
	req, rr := createSellerRequestTest(http.MethodDelete, "/sellers/4", "{}")

	r.ServeHTTP(rr, req)

	assert.Equal(t, 404, rr.Code, rr.Result())
}

func Test_delete_ok(t *testing.T) {
	r := createSellerServer()
	req, rr := createSellerRequestTest(http.MethodDelete, "/sellers/1", "{}")

	r.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code, rr.Result())
}
