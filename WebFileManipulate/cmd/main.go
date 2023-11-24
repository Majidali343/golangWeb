package main

import (
	"fmt"
	middleware "wordcount/internal/Middleware"
	filemanipulateController "wordcount/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

	router := SetupRouter()
	router.Run(":8080")
}

func Dbconnection() *gorm.DB {
	var db *gorm.DB
	var err error

	host := "localhost"
	username := "postgres"
	password := "Majid"
	dbName := "filemanipulation"

	// Construct the connection string
	connectionString := fmt.Sprintf("host=%s user=%s dbname=%s password=%s",
		host, username, dbName, password)

	// Connect to PostgreSQL
	db, err = gorm.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		
	}

	defer db.Close()
	return db
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	
	r.POST("/login", filemanipulateController.Login)
	r.POST("/register", filemanipulateController.Register)
	r.POST("/filemanipulate", middleware.AuthMiddleware(), filemanipulateController.Filemanupulate)
	r.POST("/", filemanipulateController.Details)

	return r
}
