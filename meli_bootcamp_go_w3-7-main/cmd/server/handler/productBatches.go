package handler

import (
	//"context"
	//"net/http"
	//"context"
	//"strconv"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/productBatches"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/pkg/web"
	"github.com/gin-gonic/gin"
)

type ProductBatches struct {
	service productBatches.Service
}

func NewProductBatches(p productBatches.Service) *ProductBatches {
	return &ProductBatches{
		service: p,
	}
}

func (s *ProductBatches) CreateProductBatches() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req domain.ProductBatches
		if err := c.Bind(&req); err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}
		if s.service.Exists(c.Request.Context(), req.BatchNumber) {
			web.Error(c, 409, "%s", "Ya existe un Product Batches con ese number")
			return
		}
		if req.BatchNumber == 0 {
			c.JSON(422, web.NewResponse(422, nil, "El batchNumber es requerido"))
			return
		}
		if req.CurrentQuantity == 0 {
			c.JSON(422, web.NewResponse(422, nil, "CurrentCuantity es requerido"))
			return
		}
		if req.CurrentTemperature == 0 {
			c.JSON(422, web.NewResponse(422, nil, "CurrentTemperature es requerido"))
			return
		}
		if req.DueDate == "" {
			c.JSON(422, web.NewResponse(422, nil, "DueDate es requerido"))
			return
		}
		if req.InitialQuantity == 0 {
			c.JSON(422, web.NewResponse(422, nil, "InitialQuantity es requerido"))
			return
		}
		if req.ManufacturingDate == "" {
			c.JSON(422, web.NewResponse(422, nil, "ManufacturingDate es requerido"))
			return
		}
		if req.ManufacturingHour == 0 {
			c.JSON(422, web.NewResponse(422, nil, "ManufacturingHour es requerido"))
			return
		}
		if req.MinimumTemperature == 0 {
			c.JSON(422, web.NewResponse(422, nil, "MinimumTemperature es requerido"))
			return
		}
		if req.ProductId == 0 {
			c.JSON(422, web.NewResponse(422, nil, "ProductId es requerido"))
			return
		}
		if req.SectionId == 0 {
			c.JSON(422, web.NewResponse(422, nil, "SectionId es requerido"))
			return
		}

		p, err := s.service.Save(c.Request.Context(), req)
		if err != nil {
			c.JSON(422, web.NewResponse(422, nil, err.Error()))
			return
		}
		c.JSON(201, web.NewResponse(201, p, ""))
	}
}

/* func (s *ProductBatches) GetProductBatches() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Request.URL.Query().Get("id")
		if id == "" {
			web.Error(c, 404, "%s", "No hay productos batches")
			return
		}

		idParse, errorParse := strconv.Atoi(id)
		if errorParse != nil {
			web.Error(c, 404, "%s", "no se pudo parsear")
			return
		}
		ctx := context.Background()
		serv, err := s.service.GetQuantity(ctx, idParse)
		if err != nil {
			web.Error(c, 404, "%s", "No hay productos batches con id")
			return
		}
		if serv.ID == 0 {
			web.Error(c, 404, "%s", "No hay productos batches VACIO")
			return
		}
		type ProductBatchesData struct {
			SectionID     int `json:"section_id"`
			SectionNumber int `json:"section_number"`
			ProductsCount int `json:"products_count"`
		}
		var data ProductBatchesData
		sections, errorSec := s.service.GetBySectionNumber(ctx, serv.SectionId)
	}
} */
