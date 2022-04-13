package handler

import (
	"context"
	"strconv"

	"net/http"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/carry"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/pkg/web"
	"github.com/gin-gonic/gin"
)

type Carry struct {
	carryService carry.Service
}

func NewCarry(c carry.Service) *Carry {
	return &Carry{
		carryService: c,
	}
}

// GetCarry godoc
// @Summary      Get Carry
// @Description  get Carry by ID
// @Tags         carries
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Carries ID"
// @Success      200  {object}  web.response
// @Failure      404  {object}  web.errorResponse
// @Failure      500  {object}  web.errorResponse
// @Router       /api/v1/carries/{id} [get]
func (c *Carry) Get() gin.HandlerFunc {
	return func(con *gin.Context) {
		id, err := strconv.ParseInt(con.Param("id"), 10, 64)
		if err != nil {
			web.Error(con, http.StatusBadRequest, "%s", "El id ingresado no es válido")
			return
		}

		ctx := context.Background()
		carry, err := c.carryService.Get(ctx, int(id))
		if err != nil {
			web.Error(con, http.StatusNotFound, "No se encuentra el carry con id: %d", id)
			return
		}
		web.Success(con, http.StatusOK, carry)
	}
}

// ListCarries godoc
// @Summary List carries
// @Tag Carries
// @Description get all carries
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Router /api/v1/carries [get]
func (c *Carry) GetAll() gin.HandlerFunc {
	return func(con *gin.Context) {
		ctx := context.Background()
		carries, err := c.carryService.GetAll(ctx)
		if err != nil {
			web.Error(con, http.StatusInternalServerError, "Se presentó un problema en la consulta: %s", err)
			return
		}
		web.Success(con, http.StatusOK, carries)
	}
}

// CreateCarry godoc
// @Summary Create carry
// @Tag Carries
// @Description create a new carry
// @Accept json
// @Produce json
// @Param Carry body internal.domain.newCarry true "The body to create a carry"
// @Success 201 {object} web.response
// @Router /api/v1/carries [post]
func (c *Carry) Create() gin.HandlerFunc {
	return func(con *gin.Context) {
		type NewCarry struct {
			ID          int    `json:"id"`
			CID         string `json:"cid" binding:"required"`
			CompanyName string `json:"company_name"`
			Address     string `json:"address"`
			Telephone   string `json:"telephone"`
			LocalityId  int    `json:"locality_id" binding:"required"`
			BatchNumber int    `json:"batch_number"`
		}
		reqTemp := NewCarry{}
		err := con.ShouldBind(&reqTemp)
		if err != nil {
			web.Error(con, http.StatusUnprocessableEntity, "Los campos 'cid' y 'locality_id' son requeridos para la creación del Carry. Error: %s", err)
			return
		}
		ctx := context.Background()
		if c.carryService.Exists(ctx, reqTemp.CID) {
			web.Error(con, http.StatusConflict, "%s", "Ya existe un Carry asociado al código que intenta registrar. El campo 'cid' debe ser único")
			return
		}
		req := domain.Carry(reqTemp)
		c, err := c.carryService.Save(ctx, req)
		if err != nil {
			web.Error(con, http.StatusInternalServerError, "Se presentó un problema en la creación del carry: %s", err)
		}
		web.Success(con, http.StatusCreated, c)
	}
}

// UpdateCarry godoc
// @Summary Update carry
// @Tag Carries
// @Description update a carry
// @Accept json
// @Produce json
// @Param Carry body internal.domain.newCarry true "The body to update fields of a carry"
// @Success 200 {object} web.response
// @Router /api/v1/carries [patch]
func (c *Carry) Update() gin.HandlerFunc {
	return func(con *gin.Context) {
		id, err := strconv.ParseInt(con.Param("id"), 10, 64)
		if err != nil {
			web.Error(con, http.StatusBadRequest, "%s", "El id ingresado no es válido")
			return
		}
		req := domain.Carry{}
		if err := con.ShouldBindJSON(&req); err != nil {
			web.Error(con, http.StatusBadRequest, "Error en los datos de la petición: %s", err)
			return
		}
		ctx := context.Background()
		lastCarry, err := c.carryService.Get(ctx, int(id))
		if err != nil {
			web.Error(con, http.StatusNotFound, "No se encuentra el carry con id: %d", id)
			return
		}
		newCarry := updateCarryFields(lastCarry, req, int(id))
		err = c.carryService.Update(ctx, newCarry)
		if err != nil {
			web.Error(con, http.StatusInternalServerError, "Ocurrió un error al intentar actualizar el carry: %s", err)
			return
		}
		web.Success(con, http.StatusOK, newCarry)
	}
}

// DeleteCarry godoc
// @Summary Delete a carry
// @Tag Carries
// @Description delete a carry
// @Accept json
// @Produce json
// @Param        id   path      int  true  "CarryID"
// @Success      204  {object}  web.response
// @Failure      404  {object}  web.errorResponse
// @Failure      500  {object}  web.errorResponse
// @Router       /api/v1/carries/{id} [delete]
func (c *Carry) Delete() gin.HandlerFunc {
	return func(con *gin.Context) {
		id, err := strconv.ParseInt(con.Param("id"), 10, 64)
		if err != nil {
			web.Error(con, http.StatusBadRequest, "%s", "El id ingresado no es válido")
			return
		}
		ctx := context.Background()
		carry, err := c.carryService.Get(ctx, int(id))
		if carry.ID == 0 {
			web.Error(con, http.StatusNotFound, "No se encuentra el carry con id: %d", id)
			return
		}
		err = c.carryService.Delete(ctx, int(id))
		if err != nil {
			web.Error(con, http.StatusInternalServerError, "Ocurrió un error al intentar eliminar el carry: %s", err)
			return
		}
		web.Success(con, http.StatusNoContent, "")
	}
}

// ListCarriesReport godoc
// @Summary List carries resport
// @Tag Carries report
// @Description get report of carries by locality
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Router /api/v1//localities/reportCarries [get]
func (c *Carry) GetCarriesReport() gin.HandlerFunc {
	return func(con *gin.Context) {
		if con.Query("locality_id") == "" {
			ctx := context.Background()
			carries, err := c.carryService.GetCarriesReportByLocality(ctx)
			if err != nil {
				web.Error(con, http.StatusInternalServerError, "Se presentó un problema en la consulta: %s", err)
				return
			}
			web.Success(con, http.StatusOK, carries)
		} else {
			id, err := strconv.ParseInt(con.Query("locality_id"), 10, 64)
			if err != nil {
				web.Error(con, http.StatusBadRequest, "%s", "El id ingresado no es válido")
				return
			}

			ctx := context.Background()
			carry, err := c.carryService.GetCarryReportByLocalityId(ctx, int(id))
			if err != nil {
				web.Error(con, http.StatusNotFound, "No existe un locality con el id: %d", id)
				return
			}
			web.Success(con, http.StatusOK, carry)
		}

	}

}

func updateCarryFields(lastC domain.Carry, newC domain.Carry, id int) domain.Carry {
	if newC.CID != lastC.CID && newC.CID != "" {
		lastC.CID = newC.CID
	}
	if newC.Address != lastC.Address && newC.Address != "" {
		lastC.Address = newC.Address
	}
	if newC.Telephone != lastC.Telephone && newC.Telephone != "" {
		lastC.Telephone = newC.Telephone
	}
	if newC.CompanyName != lastC.CompanyName && newC.CompanyName != "" {
		lastC.CompanyName = newC.CompanyName
	}
	if newC.LocalityId != lastC.LocalityId && newC.LocalityId != 0 {
		lastC.LocalityId = newC.LocalityId
	}
	if newC.BatchNumber != lastC.BatchNumber && newC.BatchNumber != 0 {
		lastC.BatchNumber = newC.BatchNumber
	}
	lastC.ID = id
	return lastC
}
