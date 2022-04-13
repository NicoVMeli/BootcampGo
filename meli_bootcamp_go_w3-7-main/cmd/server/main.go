package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/extlurosell/meli_bootcamp_go_w3-7/cmd/server/routes"
	"github.com/extlurosell/meli_bootcamp_go_w3-7/docs"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// @title MELI Bootcamp Sprint
// @version 1.0.0
// @description This API Handle MELI Sprint

// @lisence.name Apache 2.0
// @lisence.url https://apache.org/licenses/LICENSE-2.0
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	docs.SwaggerInfo.Host = os.Getenv("HOST")
	db, _ := sql.Open("mysql", "meli_sprint_user:Meli_Sprint#123@/melisprint")
	r := gin.Default()

	router := routes.NewRouter(r, db)
	router.MapRoutes()

	if err := r.Run(); err != nil {
		panic(err)
	}
}
