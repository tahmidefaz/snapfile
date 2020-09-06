package dbutils

import (
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"../types"
	"fmt"
)

var DB *gorm.DB

const (
    host     = "192.168.99.100"
    port     = "5432"
    user     = "snapfile"
    password = "snapfile"
    dbname   = "snapfile"
)

func DatabaseConnect() {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&types.DbModal{})

	DB = db

	fmt.Println("DB Migration Complete")
}
