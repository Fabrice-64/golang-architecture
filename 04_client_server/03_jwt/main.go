package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	godotenv "github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func goDotEnv(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func main() {
	godotenv.Load(".env")
	localhost := goDotEnv("host")
	user := goDotEnv("user")
	pwd := goDotEnv("password")
	dbname := goDotEnv("dbname")
	dsn := "host=" + localhost + " user=" + user + " password=" + pwd + " dbname=" + dbname + " port=5432 sslmode=disable"
	fmt.Println(strconv.Quote(dsn))
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqldb, err := conn.DB()
	if err != nil {
		panic(err)
	}
	err = sqldb.Ping()
	if err != nil {
		log.Fatal("database connected")
	}
	fmt.Println("CONNEXION - ENFIN")
}
