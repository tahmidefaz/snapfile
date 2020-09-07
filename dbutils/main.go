package dbutils

import (
	"fmt"
	"github.com/tahmidefaz/snapfile/misc"
	"github.com/tahmidefaz/snapfile/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var (
	host     = misc.GetEnv("DB_HOST", "0.0.0.0")
	port     = misc.GetEnv("DB_PORT", "5432")
	user     = misc.GetEnv("DB_USER", "snapfile")
	password = misc.GetEnv("DB_PASSWORD", "snapfile")
	dbname   = misc.GetEnv("DB_NAME", "snapfile")
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
