package handler

import (
	"context"
	"net/http"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/locality"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/pkg/web"
	"github.com/gin-gonic/gin"
)

type Locality struct {
	service locality.Service
}

const tokenC string = "1234"

func NewLocality(p locality.Service) *Locality {
	return &Locality{
		service: p,
	}
}
func (s *Locality) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token != tokenC {
			c.JSON(401, web.NewResponse(401, nil, "Token inválido"))
			return
		}
		var req domain.Locality
		if err := c.Bind(&req); err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}
		if s.service.Exists(c.Request.Context(), req.ZipCode) {
			web.Error(c, 409, "%s", "Ya existe un locality con ese ZipCode")
			return
		}
		if req.ZipCode == "" {
			c.JSON(422, web.NewResponse(422, nil, "El zip code es requerido"))
			return
		}
		if req.LocalityName == "" {
			c.JSON(422, web.NewResponse(422, nil, "El nombre de la  localidad es requerido"))
			return
		}
		if req.ProvinceName == "" {
			c.JSON(422, web.NewResponse(422, nil, "El nombre de la provincia es requerido"))
			return
		}
		if req.CountryName == "" {
			c.JSON(422, web.NewResponse(422, nil, "El nombre del pais es requerido"))
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
func (s *Locality) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token != tokenC {
			c.JSON(401, web.NewResponse(401, nil, "Token inválido"))
			return
		}
		zipCode := c.Request.URL.Query().Get("zip_code")
		if zipCode == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "Parametro 'zip_code' vacío"))
			return
		}
		l, err := s.service.GetByZipCode(context.Background(), zipCode)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewResponse(http.StatusInternalServerError, nil, ""))
			return
		}
		if l.ID == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, nil, ""))
			return
		}
		type LocalityData struct {
			ZipCode      string `json:"zip_code"`
			LocalityName string `json:"locality_name"`
			SellerCounts int    `json:"sellers_count"`
		}
		var data LocalityData
		sellers, err := s.service.GetSellers(context.Background(), l)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewResponse(http.StatusInternalServerError, nil, ""))
			return
		}
		data.LocalityName = l.LocalityName
		data.ZipCode = l.ZipCode
		data.SellerCounts = len(sellers)
		c.JSON(http.StatusOK, web.NewResponse(http.StatusOK, data, ""))
	}
}
