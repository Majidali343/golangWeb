package main

import (
	filemanipulateController "wordcount/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	router := SetupRouter()
	router.Run(":8080")
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/login", filemanipulateController.Login)
	r.POST("/filemanipulate", filemanipulateController.Filemanupulate)
	r.POST("/", filemanipulateController.Details)

	return r
}
