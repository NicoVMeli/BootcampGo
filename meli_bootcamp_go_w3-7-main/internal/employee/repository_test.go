package employee

import (
	"context"
	"fmt"
	"log"
	"testing"

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
	employeeTest := domain.Employee{
		CardNumberID: "123",
		FirstName:    "laytest",
		LastName:     "guetest",
		WarehouseID:  1,
	}
	employeeUpdate := domain.Employee{
		ID:           1,
		CardNumberID: "123",
		FirstName:    "hadassa",
		LastName:     "guerrero",
		WarehouseID:  2,
	}
	var expectedResultReportGetAll []domain.ReportInboundOrders
	expectedResultReportByEmployees := domain.ReportInboundOrders{}

	if _, err := db.Exec(`Truncate table employees;`); err != nil {
		log.Fatal(err)
	}

	existEmployees := repo.Exists(ctx, employeeTest.CardNumberID)
	assert.Equal(t, false, existEmployees)
	id, err := repo.Save(ctx, employeeTest)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	employees, err := repo.GetAll(ctx)
	assert.Len(t, employees, 1)
	employeeResult, err := repo.Get(ctx, 1)
	assert.NoError(t, err)
	_, err = repo.Get(ctx, 4)
	assert.Error(t, err)
	employeeTest.ID = id
	assert.Equal(t, employeeTest, employeeResult)
	err = repo.Update(ctx, employeeUpdate)
	assert.NoError(t, err)
	employeeCheckUpdate, _ := repo.Get(ctx, 1)
	assert.Equal(t, employeeUpdate, employeeCheckUpdate)
	employeeUpdate.ID = 12
	err = repo.Update(ctx, employeeUpdate)
	fmt.Println("errr----", err)
	assert.Error(t, err)
	err = repo.Delete(ctx, 1)
	assert.NoError(t, err)
	err = repo.Delete(ctx, 2)
	assert.Error(t, err)
	resultList, err := repo.GetAllReportInboundOrders()
	assert.NoError(t, err)
	assert.Equal(t, expectedResultReportGetAll, resultList)
	resultReportByEmployees, _ := repo.GetReportByEmployeeIdInboundOrders(1)
	assert.Equal(t, expectedResultReportByEmployees, resultReportByEmployees)

}
