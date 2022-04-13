package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/employee"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/pkg/web"
	"github.com/gin-gonic/gin"
)

type Employee struct {
	employeeService employee.Service
}

func NewEmployee(e employee.Service) *Employee {
	return &Employee{
		employeeService: e,
	}
}

// GetEmployee godoc
// @Summary      Get Employee
// @Description  get Employee by ID
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Employee ID"
// @Success      200  {object}  web.response
// @Failure      404  {object}  web.errorResponse
// @Failure      500  {object}  web.errorResponse
// @Router       /api/v1/employees/{id} [get]
func (e *Employee) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "Error: %s", "invalid ID")
			return
		}

		ctx := context.Background()
		employeById, err := e.employeeService.Get(ctx, int(id))
		if employeById.ID == 0 {
			web.Error(c, http.StatusNotFound, "No existe el empleado con el id %d", id)
			return
		}
		if err != nil {
			web.Error(c, http.StatusNotFound, "Error in code for: %s", err.Error())
			return
		}
		web.Success(c, http.StatusOK, employeById)

	}
}

// ListEmployees godoc
// @Summary List employees
// @Tag Employees
// @Description get all employees
// @Accept json
// @Produce json
// @Success 200 {object} web.response
// @Router /api/v1/employees [get]
func (e *Employee) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		employees, err := e.employeeService.GetAll(ctx)
		if err != nil {
			web.Error(c, http.StatusNotFound, "Error in code for: %s", err.Error())
			return
		}
		if len(employees) == 0 {
			web.Error(c, http.StatusNotFound, "%s", "No Existen Registros")
			return
		}
		web.Success(c, http.StatusOK, employees)
	}
}

func (e *Employee) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := domain.Employee{}
		err := c.Bind(&req)
		if err != nil {
			web.Error(c, http.StatusNotFound, "Error in code for: %s", err)
			return
		}
		if validateFields(req) {
			web.Error(c, http.StatusUnprocessableEntity, "%s", "Todos los campos son requeridos")
			return
		}
		ctx := context.Background()
		employee, err := e.employeeService.Save(ctx, req)
		if err != nil {
			web.Error(c, http.StatusNotFound, "%s", err)
			return
		}
		web.Success(c, http.StatusCreated, employee)
	}
}

func (e *Employee) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "error: %s", "invalid ID")
			return
		}
		req := domain.Employee{}
		if err := c.ShouldBindJSON(&req); err != nil {
			web.Error(c, http.StatusNotFound, "Error in code for: %s", err)
			return
		}
		ctx := context.Background()
		result, err := e.employeeService.Update(ctx, int(id), req)
		if err != nil {
			web.Error(c, http.StatusNotFound, "Error in code for: %s", err)
			return
		}
		web.Success(c, http.StatusOK, result)
	}
}

func (e *Employee) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "Error: %s", "invalid ID")
			return
		}
		ctx := context.Background()
		err = e.employeeService.Delete(ctx, int(id))
		if err != nil {
			web.Error(c, http.StatusNotFound, "Error in code for: %s", err)
			return
		}
		web.Success(c, http.StatusNoContent, "")
	}
}

func (e *Employee) GetReportIO() gin.HandlerFunc {
	return func(c *gin.Context) {
		idurl := c.Query("id")
		ctx := context.Background()
		if idurl == "" {
			results, err := e.employeeService.GetAllReportInboundOrders(ctx)
			/*if len(results) == 0 {
				web.Error(c, http.StatusNotFound, "%s", "No existen resultados")
				return
			}*/
			if err != nil {
				web.Error(c, http.StatusNotFound, "Error in code for: %s", err.Error())
				return
			}
			web.Success(c, http.StatusOK, results)
			return
		}

		id, err := strconv.ParseInt(idurl, 10, 64)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "Error: %s", "invalid ID")
			return
		}
		report, err := e.employeeService.GetReportByEmployeeIdInboundOrders(ctx, int(id))
		if report.ID == 0 {
			web.Error(c, http.StatusNotFound, "No existen resultados para el empleado con el id %d", id)
			return
		}
		if err != nil {
			web.Error(c, http.StatusNotFound, "Error in code for: %s", err.Error())
			return
		}
		web.Success(c, http.StatusOK, report)

	}
}

func validateFields(e domain.Employee) bool {
	if e.CardNumberID == "" {
		return true
	}
	if e.FirstName == "" {
		return true
	}
	if e.LastName == "" {
		return true
	}
	if e.WarehouseID == 0 {
		return true
	}
	return false
}
