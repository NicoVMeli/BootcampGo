package handler

import (
	"context"
	"strconv"

	"net/http"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/warehouse"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/pkg/web"
	"github.com/gin-gonic/gin"
)

type Warehouse struct {
	warehouseService warehouse.Service
}

func NewWarehouse(w warehouse.Service) *Warehouse {
	return &Warehouse{
		warehouseService: w,
	}
}

// GetWarehouse godoc
// @Summary      Get Warehouse
// @Description  get Warehouse by ID
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Warehouse ID"
// @Success      200  {object}  web.response
// @Failure      404  {object}  web.errorResponse
// @Failure      500  {object}  web.errorResponse
// @Router       /api/v1/warehouses/{id} [get]
func (w *Warehouse) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "%s", "El id ingresado no es válido")
			return
		}

		ctx := context.Background()
		warehouse, err := w.warehouseService.Get(ctx, int(id))
		if err != nil {
			web.Error(c, http.StatusNotFound, "No se encuentra el warehouse con id: %d", id)
			return
		}
		web.Success(c, http.StatusOK, warehouse)
	}
}

// ListWarehouses godoc
// @Summary List warehouses
// @Tag Warehouses
// @Description get all warehouses
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Router /api/v1/warehouses [get]
func (w *Warehouse) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		warehouses, err := w.warehouseService.GetAll(ctx)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, "Se presentó un problema en la consulta: %s", err)
			return
		}
		web.Success(c, http.StatusOK, warehouses)
	}
}

// CreateWarehouse godoc
// @Summary Create warehouse
// @Tag Warehouses
// @Description create a new warehouse
// @Accept json
// @Produce json
// @Param Warehouse body internal.domain.newWarehouse true "The body to create a warehouse"
// @Success 201 {object} web.response
// @Router /api/v1/warehouses [post]
func (w *Warehouse) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		type newWarehouse struct {
			ID                 int    `json:"id"`
			Address            string `json:"address"  binding:"required"`
			Telephone          string `json:"telephone"  binding:"required"`
			WarehouseCode      string `json:"warehouse_code"  binding:"required"`
			MinimumCapacity    int    `json:"minimum_capacity"  binding:"required"`
			MinimumTemperature int    `json:"minimum_temperature"  binding:"required"`
		}
		reqTemp := newWarehouse{}
		err := c.ShouldBind(&reqTemp)
		if err != nil {
			web.Error(c, http.StatusUnprocessableEntity, "Todos los campos son requeridos (address, telephone, warehouse_code, minimum_capacity, minimum_temperature). Error: %s", err)
			return
		}
		ctx := context.Background()
		if w.warehouseService.Exists(ctx, reqTemp.WarehouseCode) {
			web.Error(c, http.StatusConflict, "%s", "Ya existe un Warehouse asociado al código que intenta registrar. El campo 'warehouse_code' debe ser único")
			return
		}
		req := domain.Warehouse(reqTemp)
		w, err := w.warehouseService.Save(ctx, req)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, "Se presentó un problema en la creación del warehouse: %s", err)
		}
		web.Success(c, http.StatusCreated, w)
	}
}

// UpdateWarehouse godoc
// @Summary Update warehouse
// @Tag Warehouses
// @Description update a warehouse
// @Accept json
// @Produce json
// @Param Warehouse body internal.domain.newWarehouse true "The body to update fields of a warehouse"
// @Success 200 {object} web.response
// @Router /api/v1/warehouses [patch]
func (w *Warehouse) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "%s", "El id ingresado no es válido")
			return
		}
		req := domain.Warehouse{}
		if err := c.ShouldBindJSON(&req); err != nil {
			web.Error(c, http.StatusBadRequest, "Error en los datos de la petición: %s", err)
			return
		}
		ctx := context.Background()
		lastWarehouse, err := w.warehouseService.Get(ctx, int(id))
		if err != nil {
			web.Error(c, http.StatusNotFound, "No se encuentra el warehouse con id: %d", id)
			return
		}
		newWarehouse := updateFields(lastWarehouse, req, int(id))
		err = w.warehouseService.Update(ctx, newWarehouse)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, "Ocurrió un error al intentar actualizar el warehouse: %s", err)
			return
		}
		web.Success(c, http.StatusOK, newWarehouse)
	}
}

// DeleteWarehouse godoc
// @Summary Delete a warehouse
// @Tag Warehouses
// @Description delete a warehouse
// @Accept json
// @Produce json
// @Param        id   path      int  true  "Warehouse ID"
// @Success      204  {object}  web.response
// @Failure      404  {object}  web.errorResponse
// @Failure      500  {object}  web.errorResponse
// @Router       /api/v1/warehouses/{id} [delete]
func (w *Warehouse) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "%s", "El id ingresado no es válido")
			return
		}
		ctx := context.Background()
		warehouse, err := w.warehouseService.Get(ctx, int(id))
		if warehouse.ID == 0 {
			web.Error(c, http.StatusNotFound, "No se encuentra el warehouse con id: %d", id)
			return
		}
		err = w.warehouseService.Delete(ctx, int(id))
		if err != nil {
			web.Error(c, http.StatusInternalServerError, "Ocurrió un error al intentar eliminar el warehouse: %s", err)
			return
		}
		web.Success(c, http.StatusNoContent, "")
	}
}

func updateFields(lastW domain.Warehouse, newW domain.Warehouse, id int) domain.Warehouse {
	if newW.Address != lastW.Address && newW.Address != "" {
		lastW.Address = newW.Address
	}
	if newW.Telephone != lastW.Telephone && newW.Telephone != "" {
		lastW.Telephone = newW.Telephone
	}
	if newW.WarehouseCode != lastW.WarehouseCode && newW.WarehouseCode != "" {
		lastW.WarehouseCode = newW.WarehouseCode
	}
	if newW.MinimumCapacity != lastW.MinimumCapacity && newW.MinimumCapacity != 0 {
		lastW.MinimumCapacity = newW.MinimumCapacity
	}
	if newW.MinimumTemperature != lastW.MinimumTemperature && newW.MinimumTemperature != 0 {
		lastW.MinimumTemperature = newW.MinimumTemperature
	}
	lastW.ID = id
	return lastW
}
