package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/section"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/pkg/web"
	"github.com/gin-gonic/gin"
)

type Section struct {
	sectionService section.Service
}

func NewSection(s section.Service) *Section {
	return &Section{
		sectionService: s,
	}
}

func (s *Section) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		sections, err := s.sectionService.GetAll(ctx)
		if err != nil {
			web.Error(c, 404, "Error: %s", err.Error())
			return
		}
		if len(sections) == 0 {
			web.Error(c, 404, "%s", "No hay registros en la base de datos")
			return
		}
		web.Success(c, 201, sections)
	}
}

func (s *Section) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			web.Error(c, 400, "Error: %s", "The ID is not valid")
			return
		}
		ctx := context.Background()
		sectionById, err := s.sectionService.Get(ctx, int(id))
		if sectionById.ID == 0 {
			web.Error(c, 404, "there is no section with the ID %d", id)
			return
		}
		if err != nil {
			web.Error(c, 404, "Error: %s", err.Error())
			return
		}
		web.Success(c, 201, sectionById)

	}
}

func (s *Section) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req domain.Section
		if err := c.Bind(&req); err != nil {
			web.Error(c, 404, "Error: %s", err.Error())
			return
		}
		ctx := context.Background()
		fmt.Println("soy el request ", req)
		if s.sectionService.Exists(ctx, req.SectionNumber) {
			web.Error(c, 409, "%s", "Ya existe una Seccion con ese numero")
			return
		}
		if req.SectionNumber == 0 {
			web.Error(c, 422, "El campo Section Number es requerido")
			return
		}
		if req.CurrentTemperature == 0 {
			web.Error(c, 422, "El campo CurrentTemperature es requerido")
			return
		}
		if req.MinimumTemperature == 0 {
			web.Error(c, 422, "El campo MinimumTemperature es requerido")
			return
		}
		if req.CurrentCapacity == 0 {
			web.Error(c, 422, "El campo CurrentCapacity es requerido")
			return
		}
		if req.MinimumCapacity == 0 {
			web.Error(c, 422, "El campo MinimumCapacity es requerido")
			return
		}
		if req.MaximumCapacity == 0 {
			web.Error(c, 422, "El campo MaximumCapacity es requerido")
			return
		}
		sectionCreate, err := s.sectionService.Save(ctx, req)
		if err != nil {
			web.Error(c, 404, "Error: %s", err.Error())
			return
		}
		web.Success(c, 201, sectionCreate)

	}
}

func (s *Section) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			web.Error(c, 404, "Error: %s", "El ID no es valido")
			return
		}
		var req domain.Section
		if err := c.ShouldBindJSON(&req); err != nil {
			web.Error(c, 404, "Error in code for: %s", err)
			return
		}
		ctx := context.Background()
		sectionById, err := s.sectionService.Get(ctx, int(id))
		if err != nil {
			web.Error(c, 400, "Error: %s", err.Error())
			return
		}
		if sectionById.ID == 0 {
			web.Error(c, 404, "Error: %s", "The ID is not valid")
			return
		}
		if req.SectionNumber == 0 {
			web.Error(c, 400, "El campo Section Number es requerido")
			return
		}
		if req.CurrentTemperature == 0 {
			web.Error(c, 400, "El campo CurrentTemperature es requerido")
			return
		}
		if req.MinimumTemperature == 0 {
			web.Error(c, 400, "El campo MinimumTemperature es requerido")
			return
		}
		if req.CurrentCapacity == 0 {
			web.Error(c, 400, "El campo CurrentCapacity es requerido")
			return
		}
		if req.MinimumCapacity == 0 {
			web.Error(c, 400, "El campo MinimumCapacity es requerido")
			return
		}
		if req.MaximumCapacity == 0 {
			web.Error(c, 400, "El campo MaximumCapacity es requerido")
			return
		}
		req.ID = int(id)
		sectionUpdate, err := s.sectionService.Update(ctx, req)
		if err != nil {
			web.Error(c, 400, "Error: %s", err.Error())
			return
		}
		web.Success(c, 200, sectionUpdate)
	}
}

func (s *Section) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			web.Error(c, 404, "Error: %s", "Invalido ID")
			return
		}
		ctx := context.Background()
		sectionById, err := s.sectionService.Get(ctx, int(id))
		if sectionById.ID == 0 {
			web.Error(c, 404, "No existe seccion con ese ID")
			return
		}
		if err != nil {
			web.Error(c, 400, "Error: %s", err.Error())
			return
		}
		sectionDelete := s.sectionService.Delete(ctx, int(id))
		if sectionDelete != nil {
			web.Error(c, 404, "No se pudo remover la seccion")
			return
		}
		web.Success(c, 204, "La seccion fue removida")
	}
}
