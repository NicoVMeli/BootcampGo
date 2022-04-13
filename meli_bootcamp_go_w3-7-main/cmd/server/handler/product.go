package handler

import (
	"errors"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/product"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/pkg/web"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Product struct {
	productService product.Service
}

func NewProduct(productService product.Service) *Product {
	return &Product{
		productService: productService,
	}
}

func (p *Product) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := p.productService.GetAll(c)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, "error: %s", err.Error())
			return
		}
		web.Success(c, http.StatusOK, products)
	}
}

func (p *Product) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := c.Params.Get("id")
		if !ok {
			web.Error(c, http.StatusInternalServerError, "error: unable to retrieve 'id' param from URL")
			return
		}
		idInt, err := strconv.ParseInt(id, 10, 0)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "error: provided id '%s' is not an integer", id)
			return
		}
		product, err := p.productService.Get(c, int(idInt))
		if err != nil {
			web.Error(c, http.StatusNotFound, "error: %s", err.Error())
			return
		}
		web.Success(c, http.StatusOK, product)
	}
}

func (p *Product) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productToSave domain.Product
		// Asociar el contenido del body a los campos de la estructura Product
		if err := c.ShouldBindJSON(&productToSave); err != nil {
			web.Error(c, http.StatusBadRequest, "error: %s", err.Error())
			return
		}
		if err := validateProductStruct(c, productToSave); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, "error: %s", err.Error())
			return
		}
		savedProduct, err := p.productService.Save(c, productToSave)
		if err != nil {
			if err == product.ErrAlreadyExists {
				web.Error(c, http.StatusConflict, "error: %s", err.Error())
				return
			}
			web.Error(c, http.StatusBadRequest, "error: %s", err.Error())
			return
		}
		web.Success(c, http.StatusCreated, savedProduct)
	}
}

func (p *Product) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := c.Params.Get("id")
		if !ok {
			web.Error(c, http.StatusInternalServerError, "error: unable to retrieve 'id' param from URL")
			return
		}
		idInt, err := strconv.ParseInt(id, 10, 0)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "error: provided id '%s' is not an integer", id)
			return
		}

		var productToUpdate domain.Product
		// Asociar el contenido del body a los campos de la estructura Product
		if err := c.ShouldBindJSON(&productToUpdate); err != nil {
			web.Error(c, http.StatusBadRequest, "error: %s", err.Error())
			return
		}
		updatedProduct, err := p.productService.Update(c, int(idInt), productToUpdate)
		if err != nil {
			if err.Error() == "product not found" {
				web.Error(c, http.StatusNotFound, "error: %s", err.Error())
				return
			}
			web.Error(c, http.StatusInternalServerError, "error: %s", err.Error())
			return
		}
		web.Success(c, http.StatusOK, updatedProduct)
	}
}

func (p *Product) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := c.Params.Get("id")
		if !ok {
			web.Error(c, http.StatusInternalServerError, "error: unable to retrieve 'id' param from URL")
			return
		}
		idInt, err := strconv.ParseInt(id, 10, 0)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "error: provided id '%s' is not an integer", id)
			return
		}
		err = p.productService.Delete(c, int(idInt))
		if err != nil {
			web.Error(c, http.StatusNotFound, "error: %s", err.Error())
			return
		}
		web.Success(c, http.StatusNoContent, nil)
	}
}

func (p *Product) GetRecordReportsByProductId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Define Response Data struct
		type responseData struct {
			ProductId           int    `json:"product_id"`
			Description         string `json:"description"`
			ProductRecordsCount int    `json:"records_count"`
		}
		// Get query param "id" from URL
		id, ok := c.GetQuery("id")
		if !ok {
			web.Error(c, http.StatusInternalServerError, "error: unable to retrieve 'id' query param from URL")
			return
		}
		// Validate that id is an integer
		idInt, err := strconv.ParseInt(id, 10, 0)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "error: provided id '%s' is not an integer", id)
			return
		}
		// Get Product Record reports
		description, productRecordsCount, err := p.productService.GetRecordReportsByProductId(c, int(idInt))
		if err != nil {
			web.Error(c, http.StatusNotFound, "error: %s", err.Error())
			return
		}
		// Set the response data
		var response = responseData{
			ProductId:           int(idInt),
			Description:         description,
			ProductRecordsCount: productRecordsCount,
		}
		web.Success(c, http.StatusOK, response)
	}
}

func validateProductStruct(ctx *gin.Context, product domain.Product) error {
	if product.ProductCode == "" {
		return errors.New("product code is required")
	}
	if product.Description == "" {
		return errors.New("product description is required")
	}
	if product.Width == 0 {
		return errors.New("product width is required")
	}
	if product.Height == 0 {
		return errors.New("product height is required")
	}
	if product.Length == 0 {
		return errors.New("product length is required")
	}
	if product.Netweight == 0 {
		return errors.New("product netweight is required")
	}
	if product.ExpirationRate == 0 {
		return errors.New("product expiration rate is required")
	}
	if product.RecomFreezTemp == 0 {
		return errors.New("product recommended freezing temperature is required")
	}
	if product.FreezingRate == 0 {
		return errors.New("product freezing rate is required")
	}
	if product.ProductTypeID == 0 {
		return errors.New("product type id is required")
	}
	return nil
}
