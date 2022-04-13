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

func createSectionServer() *gin.Engine {
	sectionService := NewSectionServiceMock()
	section := NewSection(sectionService)
	r := gin.Default()
	sectionsGroup := r.Group("/sections")
	{
		sectionsGroup.GET("/", section.GetAll())
		sectionsGroup.GET("/:id", section.Get())
		sectionsGroup.POST("/", section.Create())
		sectionsGroup.PATCH("/:id", section.Update())
		sectionsGroup.DELETE("/:id", section.Delete())
	}
	return r
}
func createSectionRequestTest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

type SectionServiceMock struct {
	Data []byte
}

func NewSectionServiceMock() *SectionServiceMock {
	datosDePrueba := []domain.Section{
		{
			ID:                 1,
			SectionNumber:      41,
			CurrentTemperature: 3,
			MinimumTemperature: 3,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    3,
			WarehouseID:        3,
			ProductTypeID:      3,
		},
		{
			ID:                 2,
			SectionNumber:      42,
			CurrentTemperature: 3,
			MinimumTemperature: 2,
			CurrentCapacity:    2,
			MinimumCapacity:    2,
			MaximumCapacity:    2,
			WarehouseID:        2,
			ProductTypeID:      2,
		},
		{
			ID:                 3,
			SectionNumber:      43,
			CurrentTemperature: 32,
			MinimumTemperature: 22,
			CurrentCapacity:    22,
			MinimumCapacity:    22,
			MaximumCapacity:    22,
			WarehouseID:        22,
			ProductTypeID:      22,
		},
	}
	datos, _ := json.Marshal(datosDePrueba)
	return &SectionServiceMock{
		Data: datos,
	}
}

func (s *SectionServiceMock) GetAll(ctx context.Context) ([]domain.Section, error) {
	var datos []domain.Section
	err := json.Unmarshal(s.Data, &datos)
	if err != nil {
		return []domain.Section{}, err
	}
	return datos, nil
}
func (s *SectionServiceMock) Get(ctx context.Context, id int) (domain.Section, error) {
	var datos []domain.Section
	err := json.Unmarshal(s.Data, &datos)
	if err != nil {
		return domain.Section{}, err
	}
	for _, section := range datos {
		if section.ID == id {
			return section, nil
		}
	}
	return domain.Section{}, fmt.Errorf("section not found")
}
func (s *SectionServiceMock) Save(ctx context.Context, se domain.Section) (domain.Section, error) {
	var datos []domain.Section
	err := json.Unmarshal(s.Data, &datos)
	if err != nil {
		return domain.Section{}, err
	}
	for _, sect := range datos {
		if sect.SectionNumber == se.SectionNumber {
			return domain.Section{}, fmt.Errorf("section already exists")
		}
	}
	se.ID = datos[len(datos)-1].ID
	datos = append(datos, se)
	return se, nil
}
func (s *SectionServiceMock) Update(ctx context.Context, se domain.Section) (domain.Section, error) {
	var datos []domain.Section
	err := json.Unmarshal(s.Data, &datos)
	if err != nil {
		return domain.Section{}, err
	}
	for _, sec := range datos {
		if sec.ID == se.ID {
			sec.SectionNumber = se.SectionNumber
			sec.MinimumCapacity = se.MinimumCapacity
			sec.MinimumTemperature = se.MinimumTemperature
			sec.MaximumCapacity = se.MaximumCapacity
			sec.CurrentCapacity = se.CurrentCapacity
			sec.CurrentTemperature = se.CurrentTemperature
			sec.WarehouseID = se.WarehouseID
			sec.ProductTypeID = se.ProductTypeID
			return sec, nil
		}
	}
	return domain.Section{}, fmt.Errorf("Section does not exist")
}
func (s *SectionServiceMock) Delete(ctx context.Context, id int) error {
	var datos []domain.Section
	err := json.Unmarshal(s.Data, &datos)
	if err != nil {
		return err
	}
	for _, section := range datos {
		if section.ID == id {
			return nil
		}
	}
	return fmt.Errorf("Section not found")
}

//Sacarlo del service, asi no rompe por todos lados.
func (s *SectionServiceMock) Exists(ctx context.Context, sectionNumber int) bool {
	return false
}

//trabado
func Test_create_ok_section(t *testing.T) {
	r := createSectionServer()
	body := `{"section_number": 89,"current_temperature": 3,"minimum_temperature": 3,"current_capacity": 1,"minimum_capacity": 1,"maximum_capacity": 3,"warehouse_id": 3,"product_type_id": 3}`
	req, rr := createSectionRequestTest(http.MethodPost, "/sections/", body)
	r.ServeHTTP(rr, req)
	assert.Equal(t, 201, rr.Code, rr.Result())
}

func Test_create_fail_section(t *testing.T) {
	r := createSectionServer()
	body := `
        {
			"section_number": 5,
			"current_temperature": 3,
			"minimum_temperature": 3,
			"current_capacity": 1,
			"minimum_capacity": 1,
        }`
	req, rr := createSectionRequestTest(http.MethodPost, "/sections/", body)
	r.ServeHTTP(rr, req)
	assert.Equal(t, 400, rr.Code, rr.Result())
}

func Test_create_conflict_section(t *testing.T) {
	r := createSectionServer()
	body := `
        {
			"section_number": 41,
			"current_temperature": 3,
			"minimum_temperature": 3,
			"current_capacity": 1,
			"minimum_capacity": 1,
			"maximum_capacity": 3,
			"warehouse_id": 3,
			"product_type_id": 3
        }`
	req, rr := createSectionRequestTest(http.MethodPost, "/sections/", body)
	r.ServeHTTP(rr, req)
	assert.Equal(t, 404, rr.Code, rr.Result())
}

func Test_find_all_section(t *testing.T) {
	r := createSectionServer()
	req, rr := createSectionRequestTest(http.MethodGet, "/sections/", "{}")
	r.ServeHTTP(rr, req)
	assert.Equal(t, 201, rr.Code, rr.Result())
}

func Test_find_by_id_non_existent_section(t *testing.T) {
	r := createSectionServer()
	req, rr := createSectionRequestTest(http.MethodGet, "/sections/17", "{}")
	r.ServeHTTP(rr, req)
	assert.Equal(t, 404, rr.Code, rr.Result())
}

func Test_find_by_id_existent_section(t *testing.T) {
	r := createSectionServer()
	req, rr := createSectionRequestTest(http.MethodGet, "/sections/1", "{}")
	r.ServeHTTP(rr, req)
	assert.Equal(t, 201, rr.Code, rr.Result())
}

func Test_update_ok_section(t *testing.T) {
	r := createSectionServer()
	body := `
        {
			"section_number": 41,
			"current_temperature": 3,
			"minimum_temperature": 3,
			"current_capacity": 1,
			"minimum_capacity": 1,
			"maximum_capacity": 3,
			"warehouse_id": 3,
			"product_type_id": 3
        }`
	req, rr := createSectionRequestTest(http.MethodPatch, "/sections/1", body)
	r.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code)
}

func Test_update_non_existent_section(t *testing.T) {
	r := createSectionServer()
	body := `
        {
			"section_number": 41,
			"current_temperature": 3,
			"minimum_temperature": 3,
			"current_capacity": 1,
			"minimum_capacity": 1,
			"maximum_capacity": 3,
			"warehouse_id": 3,
			"product_type_id": 3
        }`
	req, rr := createSectionRequestTest(http.MethodPatch, "/sections/15", body)
	r.ServeHTTP(rr, req)
	assert.Equal(t, 400, rr.Code)
}

func Test_delete_non_existent_section(t *testing.T) {
	r := createSectionServer()
	req, rr := createSectionRequestTest(http.MethodDelete, "/sections/14", "{}")
	r.ServeHTTP(rr, req)
	assert.Equal(t, 404, rr.Code, rr.Result())
}

func Test_delete_ok_section(t *testing.T) {
	r := createSectionServer()
	req, rr := createSectionRequestTest(http.MethodDelete, "/sections/1", "{}")
	r.ServeHTTP(rr, req)
	assert.Equal(t, 204, rr.Code, rr.Result())
}
