package main

import (
	middleware "wordcount/internal/Middleware"
	filemanipulateController "wordcount/internal/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"

     "github.com/gin-contrib/cors"
)

func main() {

	router := SetupRouter()
	config := cors.Config{
		AllowOrigins:     []string{"https://editor.swagger.io"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization"},
		AllowCredentials: true,
	}

	router.Use(cors.New(config))
	router.Run(":8080")
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/login", filemanipulateController.Login)
	r.POST("/admin", filemanipulateController.Adminlogin)
	r.POST("/register", filemanipulateController.Register)
	r.POST("/filemanipulate", middleware.AuthMiddleware(), filemanipulateController.Filemanupulate)
	r.GET("/UserFileStatics", middleware.AuthMiddleware(), filemanipulateController.UserFileStatics)
	r.GET("/Admingetresults", middleware.AuthMiddleware(), filemanipulateController.Admingetresults)

	r.GET("/", filemanipulateController.Details)

	return r
}
