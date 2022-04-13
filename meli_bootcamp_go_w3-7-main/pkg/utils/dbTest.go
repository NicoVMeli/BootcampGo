package utils

import (
	"database/sql"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func init() {
	txdb.Register("txdb", "mysql", "meli_sprint_user:Meli_Sprint#123@/melisprint")
}

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("txdb", uuid.New().String())
	if err == nil {
		return db, db.Ping()
	}

	return db, err
}
