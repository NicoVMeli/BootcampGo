package handler

import (
	"fmt"
	"strconv"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/seller"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/pkg/web"
	"github.com/gin-gonic/gin"
)

type Seller struct {
	service seller.Service
}

const tokenCompare string = "1234"

func NewSeller(p seller.Service) *Seller {
	return &Seller{
		service: p, //Prueba
	}
}

func (s *Seller) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")

		if token != tokenCompare {
			c.JSON(401, web.NewResponse(401, nil, "Token inválido"))
			return
		}

		p, err := s.service.GetAll(c.Request.Context())
		if err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}

		if len(p) == 0 {
			c.JSON(404, web.NewResponse(404, nil, "No hay sellers"))
			return
		}
		c.JSON(200, web.NewResponse(200, p, ""))

	}
}

func (s *Seller) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")

		if token != tokenCompare {
			c.JSON(401, web.NewResponse(401, nil, "Token inválido"))
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid ID"})
			return
		}

		p, err := s.service.Get(c.Request.Context(), int(id))

		if err != nil {
			c.JSON(401, web.NewResponse(401, nil, "ID inválido"))
			return
		}

		c.JSON(200, p)

	}
}

func (s *Seller) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")

		if token != tokenCompare {
			c.JSON(401, web.NewResponse(401, nil, "Token inválido"))
			return
		}

		var req domain.Seller

		if err := c.Bind(&req); err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}

		if s.service.Exists(c.Request.Context(), req.CID) {
			web.Error(c, 409, "%s", "Ya existe un seller con ese CID")
			return
		}

		if req.CID == 0 {
			c.JSON(422, web.NewResponse(422, nil, "El cid es requerido"))
			return
		}

		if req.Address == "" {
			c.JSON(422, web.NewResponse(422, nil, "El domicilio es requerido"))
			return
		}

		if req.CompanyName == "" {
			c.JSON(422, web.NewResponse(422, nil, "El nombre de la compañia es requerido"))
			return
		}

		if req.Telephone == "" {
			c.JSON(422, web.NewResponse(422, nil, "El telefono es requerido"))
			return
		}

		if req.LocalitiesId == 0 {
			c.JSON(422, web.NewResponse(422, nil, "El localities_Id es requerido"))
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

func (s *Seller) Update() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Request.Header.Get("token")
		if token != tokenCompare {
			c.JSON(401, gin.H{"error": "token inválido"})
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(404, gin.H{"error": "invalid ID"})
			return
		}

		var req domain.Seller

		if err := c.Bind(&req); err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		if id != int64(req.ID) {
			c.JSON(400, gin.H{"error": "Los id no coinciden. Por favor ingrese uno correctamente"})
		}
		if req.CID == 0 {
			c.JSON(400, gin.H{"error": "El Cid es requerido"})
			return
		}

		if req.Address == "" {
			c.JSON(400, gin.H{"error": "El domicilio es requerido"})
			return
		}

		if req.CompanyName == "" {
			c.JSON(400, gin.H{"error": "El nombre de la compañia es requerido"})
			return
		}
		if req.Telephone == "" {
			c.JSON(400, gin.H{"error": "El telefono es requerido"})
			return
		}

		p, err := s.service.Update(c.Request.Context(), req)
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, p)
	}
}

func (s *Seller) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token != tokenCompare {
			c.JSON(401, gin.H{"error": "token inválido"})
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid ID"})
			return
		}

		err = s.service.Delete(c.Request.Context(), int(id))
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": fmt.Sprintf("El seller %d ha sido eliminado", id)})

	}
}
