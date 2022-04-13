package inboudOrders

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/domain"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	db, err := utils.InitDB()
	assert.NoError(t, err)
	defer db.Close()
	repo := NewRepository(db)

	ctx := context.TODO()

	orderDate, _ := time.Parse("2006-01-02", "2022-01-09")
	inboudOrders := domain.InboudOrders{
		ID:             1,
		OrderDate:      orderDate,
		OrderNumber:    "123",
		EmployeeId:     1,
		ProductBatchId: 2,
		WarehouseID:    3,
	}
	if _, err := db.Exec(`Truncate table inbound_orders;`); err != nil {
		log.Fatal(err)
	}

	id, err := repo.Save(ctx, inboudOrders)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	existsEmployee := repo.ExistsEmployee(ctx, 1)
	assert.Equal(t, false, existsEmployee)
}
