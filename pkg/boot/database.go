package boot

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	log.Println("Connecting to database...")
	start := time.Now()

	host := Environment.GetEnv("POSTGRES_HOST")
	port := Environment.GetEnv("POSTGRES_PORT")
	user := Environment.GetEnv("POSTGRES_USER")
	pass := Environment.GetEnv("POSTGRES_PASS")
	db := Environment.GetEnv("POSTGRES_DB")

	uri := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, db)

	DB, err = gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	duration := time.Since(start)
	log.Printf("Connected to database in %s\n", duration)
}
