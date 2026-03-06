package initializers

import (
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB(){
	var err error
	dsn := "postgresql://postgres:[ILovePostgre]@db.uvojlicbqfvfaidqhdmh.supabase.co:5432/postgres"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}