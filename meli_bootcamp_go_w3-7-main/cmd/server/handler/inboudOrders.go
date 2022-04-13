package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/inboudOrders"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/pkg/web"
	"github.com/gin-gonic/gin"
)

type InboudOrders struct {
	inboudOrdersService inboudOrders.Service
}

func NewInboudOrders(i inboudOrders.Service) *InboudOrders {
	return &InboudOrders{
		inboudOrdersService: i,
	}
}

func (i *InboudOrders) Create() gin.HandlerFunc {
	type request struct {
		ID             int    `json:"id"`
		OrderDate      string `json:"order_date" binding:"required"`
		OrderNumber    string `json:"order_number" binding:"required"`
		EmployeeId     int    `json:"employee_id" binding:"required"`
		ProductBatchId int    `json:"product_batch_id" binding:"required"`
		WarehouseID    int    `json:"warehouse_id" binding:"required"`
	}
	return func(c *gin.Context) {
		req := request{}
		ctx := context.Background()
		err := c.ShouldBind(&req)
		if err != nil {
			switch {
			case strings.Contains(string(err.Error()), "required"):
				web.Error(c, http.StatusUnprocessableEntity, "%s", "Todos los campos son requeridos")
				return
			default:
				web.Error(c, http.StatusNotFound, "Error in code for: %s", err)
				return
			}
		}

		orderDate, _ := time.Parse("2006-01-02", req.OrderDate)
		newInboudOrders := domain.InboudOrders{
			OrderDate:      orderDate,
			OrderNumber:    req.OrderNumber,
			EmployeeId:     req.EmployeeId,
			ProductBatchId: req.ProductBatchId,
			WarehouseID:    req.WarehouseID,
		}
		inboudOrdersResul, err := i.inboudOrdersService.Save(ctx, newInboudOrders)
		if err != nil {
			switch {
			case errors.Is(err, inboudOrders.ErrorEmployeesNoExist):
				web.Error(c, http.StatusConflict, "%s", err)
				return
			default:
				web.Error(c, http.StatusNotFound, "%s", err)
				return
			}
		}
		web.Success(c, http.StatusCreated, inboudOrdersResul)
	}
}
