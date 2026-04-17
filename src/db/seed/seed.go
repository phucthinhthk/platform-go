package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID   uint64 `gorm:"primaryKey"`
	Name string
}

func main() {
	// Nạp env từ file .env.local trong src/
	godotenv.Load(".env.local")

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASS")
		host := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")

		// Fallback nếu chạy bên ngoài Docker hoặc thiếu biến lẻ
		if user == "" { user = "root" }
		if pass == "" { pass = "root" }
		if host == "" { host = "localhost" } // Nếu chạy từ host thì trỏ localhost
		if dbPort == "" { dbPort = "3309" }
		if dbName == "" { dbName = "platform_db" }

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, dbPort, dbName)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	fmt.Println("Seeding users...")
	users := []User{
		{Name: "User 1"},
		{Name: "User 2"},
		{Name: "User 3"},
	}

	for _, u := range users {
		if err := db.Create(&u).Error; err != nil {
			log.Printf("Could not seed user %s: %v", u.Name, err)
		} else {
			fmt.Printf("Seeded user: %s\n", u.Name)
		}
	}

	fmt.Println("Seeding completed successfully.")
}
