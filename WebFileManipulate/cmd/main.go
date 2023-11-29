package main

import (
	middleware "wordcount/internal/Middleware"
	filemanipulateController "wordcount/internal/handlers"

	tokens "wordcount/internal/Middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

	router := SetupRouter()
	router.Run(":8080")
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/login", filemanipulateController.Login)
	r.POST("/admin", filemanipulateController.Adminlogin)
	r.POST("/register", filemanipulateController.Register)

	// Endpoint for token refresh
	r.POST("/refresh", tokens.RefreshTokenHandler)

	r.POST("/filemanipulate", middleware.AuthMiddleware(), filemanipulateController.Filemanupulate)
	r.GET("/UserFileStatics", middleware.AuthMiddleware(), filemanipulateController.UserFileStatics)
	r.GET("/Admingetresults", middleware.AuthMiddleware(), filemanipulateController.Admingetresults)

	r.GET("/", filemanipulateController.Details)

	return r
}
