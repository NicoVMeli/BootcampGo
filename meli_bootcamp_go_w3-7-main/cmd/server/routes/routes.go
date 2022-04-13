package routes

import (
	"database/sql"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/product_record"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/cmd/server/handler"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/carry"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/employee"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/inboudOrders"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/locality"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/product"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/section"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/seller"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/warehouse"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/internal/productBatches"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Router interface {
	MapRoutes()
}

type router struct {
	r  *gin.Engine
	rg *gin.RouterGroup
	db *sql.DB
}

func NewRouter(r *gin.Engine, db *sql.DB) Router {
	return &router{r: r, db: db}
}

func (r *router) MapRoutes() {
	r.buildSwaggerRoutes()
	r.setGroup()
	r.buildSellerRoutes()
	r.buildProductRoutes()
	r.buildSectionRoutes()
	r.buildWarehouseRoutes()
	r.buildEmployeeRoutes()
	r.buildInboudOrdersRoutes()
	r.buildBuyerRoutes()
	r.buildProductRecordRoutes()
	r.buildCarryRoutes()
	r.buildLocalityRoutes()
	r.buildProductBatchesRoutes()
}

func (r *router) buildSwaggerRoutes() {
	r.r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (r *router) setGroup() {
	r.rg = r.r.Group("/api/v1")
}

func (r *router) buildSellerRoutes() {
	// Example
	repo := seller.NewRepository(r.db)
	service := seller.NewService(repo)
	handler := handler.NewSeller(service)
	r.rg.GET("/sellers", handler.GetAll())
	r.rg.GET("/sellers/:id", handler.Get())
	r.rg.POST("/sellers", handler.Create())
	r.rg.DELETE("/sellers/:id", handler.Delete())
	r.rg.PATCH("/sellers/:id", handler.Update())
}

func (r *router) buildProductRoutes() {

	// Inyección de dependencias para capas de Producto
	var productRepository product.Repository = product.NewRepository(r.db)
	var productService product.Service = product.NewService(productRepository)
	var productController *handler.Product = handler.NewProduct(productService)

	// Captura las solicitudes para Producto
	productsRouter := r.rg.Group("/products")
	{
		productsRouter.POST("/", productController.Create())
		productsRouter.GET("/:id", productController.Get())
		productsRouter.GET("/", productController.GetAll())
		productsRouter.PATCH("/:id", productController.Update())
		productsRouter.DELETE("/:id", productController.Delete())
		productsRouter.GET("/reportRecords", productController.GetRecordReportsByProductId())
	}
}

func (r *router) buildSectionRoutes() {
	repository := section.NewRepository(r.db)
	service := section.NewService(repository)
	handler := handler.NewSection(service)
	r.rg.GET("/sections", handler.GetAll())
	r.rg.GET("/sections/:id", handler.Get())
	r.rg.POST("/sections", handler.Create())
	r.rg.PATCH("/sections/:id", handler.Update())
	r.rg.DELETE("/sections/:id", handler.Delete())
}

func (r *router) buildWarehouseRoutes() {
	warehouseRepo := warehouse.NewRepository(r.db)
	warehouseService := warehouse.NewService(warehouseRepo)
	warehouseHandler := handler.NewWarehouse(warehouseService)
	r.rg.GET("/warehouses", warehouseHandler.GetAll())
	r.rg.GET("/warehouses/:id", warehouseHandler.Get())
	r.rg.POST("/warehouses", warehouseHandler.Create())
	r.rg.PATCH("/warehouses/:id", warehouseHandler.Update())
	r.rg.DELETE("/warehouses/:id", warehouseHandler.Delete())
}

func (r *router) buildEmployeeRoutes() {
	repo := employee.NewRepository(r.db)
	service := employee.NewService(repo)
	handler := handler.NewEmployee(service)
	employeesRouter := r.rg.Group("/employees")
	{
		employeesRouter.GET("/", handler.GetAll())
		employeesRouter.GET("/:id", handler.Get())
		employeesRouter.GET("/reportinboundorders", handler.GetReportIO())
		employeesRouter.POST("/", handler.Create())
		employeesRouter.PATCH("/:id", handler.Update())
		employeesRouter.DELETE("/:id", handler.Delete())
	}
}

func (r *router) buildInboudOrdersRoutes() {
	repo := inboudOrders.NewRepository(r.db)
	service := inboudOrders.NewService(repo)
	handler := handler.NewInboudOrders(service)
	r.rg.POST("/inboundOrders", handler.Create())
}

func (r *router) buildProductRecordRoutes() {

	// Inyección de dependencias para capas de Product Record
	var productRepository product.Repository = product.NewRepository(r.db)
	var productService product.Service = product.NewService(productRepository)
	var productRecordRepository product_record.Repository = product_record.NewRepository(r.db)
	var productRecordService product_record.Service = product_record.NewService(productRecordRepository, productService)
	var productRecordController *handler.ProductRecord = handler.NewProductRecord(productRecordService)

	// Captura las solicitudes para Product Record
	productsRouter := r.rg.Group("/productRecords")
	{
		productsRouter.POST("/", productRecordController.Create())
	}
}

func (r *router) buildCarryRoutes() {
	carryRepo := carry.NewRepository(r.db)
	carryService := carry.NewService(carryRepo)
	carryHandler := handler.NewCarry(carryService)
	r.rg.GET("/carries", carryHandler.GetAll())
	r.rg.GET("/carries/:id", carryHandler.Get())
	r.rg.POST("/carries", carryHandler.Create())
	r.rg.PATCH("/carries/:id", carryHandler.Update())
	r.rg.DELETE("/carries/:id", carryHandler.Delete())

	r.rg.GET("/localities/reportCarries", carryHandler.GetCarriesReport())
}

func (r *router) buildLocalityRoutes() {
	repo := locality.NewRepository(r.db)
	service := locality.NewService(repo)
	handler := handler.NewLocality(service)
	r.rg.POST("/localities", handler.Create())
	r.rg.GET("/localities/reportSellers", handler.Get())
}

func (r *router) buildProductBatchesRoutes() {
	repo := productBatches.NewRepository(r.db)
	service := productBatches.NewService(repo)
	handler := handler.NewProductBatches(service)
	r.rg.POST("/productBatches", handler.CreateProductBatches())
	//r.rg.GET("/sections/reportProducts", handler.GetProductBatches())
}

func (r *router) buildBuyerRoutes() {}
