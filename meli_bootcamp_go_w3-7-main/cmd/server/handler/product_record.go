package handler

import (
	"errors"
	"fmt"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/product_record"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/pkg/web"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type requestData struct {
	LastUpdateDate string  `json:"last_update_date" binding:"required"`
	PurchasePrice  float64 `json:"purchase_price" binding:"required"`
	SalePrice      float64 `json:"sale_price" binding:"required"`
	ProductID      int     `json:"product_id" binding:"required"`
}

type request struct {
	Data requestData `json:"data"`
}

type responseData struct {
	ID             int     `json:"id"`
	LastUpdateDate string  `json:"last_update_date"`
	PurchasePrice  float64 `json:"purchase_price"`
	SalePrice      float64 `json:"sale_price"`
	ProductID      int     `json:"product_id"`
}

type ProductRecord struct {
	productRecordService product_record.Service
}

func NewProductRecord(productRecordService product_record.Service) *ProductRecord {
	return &ProductRecord{
		productRecordService: productRecordService,
	}
}

func (pr *ProductRecord) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestContent request
		// Asociar el contenido del body a los campos de la estructura Product Record
		if err := ctx.ShouldBindJSON(&requestContent); err != nil {
			web.Error(ctx, http.StatusBadRequest, "error: %s", err.Error())
			return
		}
		if err := validateRequestData(ctx, requestContent.Data); err != nil {
			web.Error(ctx, http.StatusUnprocessableEntity, "error: %s", err.Error())
			return
		}

		// Parse last_update_date from string to time.Time
		parsedLastUpdateDate, err := time.Parse("2006-01-02", requestContent.Data.LastUpdateDate)
		if err != nil {
			web.Error(ctx, http.StatusUnprocessableEntity, "error: %s", err.Error())
			return
		}
		// Create Product Record
		productRecordToSave := domain.ProductRecord{
			ID:             0,
			LastUpdateDate: &parsedLastUpdateDate,
			PurchasePrice:  requestContent.Data.PurchasePrice,
			SalePrice:      requestContent.Data.SalePrice,
			ProductID:      requestContent.Data.ProductID,
		}
		// Save Product Record
		savedProductRecord, err := pr.productRecordService.Save(ctx, productRecordToSave)
		// Handle returned errors
		if err != nil {
			if errors.Is(err, product_record.ErrProductNotFound) {
				web.Error(ctx, http.StatusConflict, "error: %s", err.Error())
				return
			}
			web.Error(ctx, http.StatusBadRequest, "error: %s", err.Error())
			return
		}
		// Format response data
		lastUpdateDateFormattedToString := savedProductRecord.LastUpdateDate.Format("2006-01-02")
		//fmt.Printf("%+v", lastUpdateDateFormattedToString)
		var responseData = responseData{
			ID:             savedProductRecord.ID,
			LastUpdateDate: lastUpdateDateFormattedToString,
			PurchasePrice:  savedProductRecord.PurchasePrice,
			SalePrice:      savedProductRecord.SalePrice,
			ProductID:      savedProductRecord.ProductID,
		}
		fmt.Printf("%+v", responseData)
		web.Success(ctx, http.StatusCreated, responseData)
	}
}

func validateRequestData(ctx *gin.Context, requestData requestData) error {
	if requestData.LastUpdateDate == "" {
		return errors.New("last update date is required")
	}
	if requestData.PurchasePrice == 0 {
		return errors.New("purchase price is required")
	}
	if requestData.SalePrice == 0 {
		return errors.New("sale price is required")
	}
	if requestData.ProductID == 0 {
		return errors.New("product id is required")
	}
	return nil
}
