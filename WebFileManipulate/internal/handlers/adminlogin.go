package filemanipulate

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	tokens "wordcount/internal/Middleware"
)

type Admin struct {
	Email    string `json:"email" `
	Password string `json:"password"`
	// Add other fields as needed
}

func Adminlogin(c *gin.Context) {

	var admin Admin

	if err := c.BindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if admin.Email == "admin@gmail.com" && admin.Password == "12345" {

		accessToken, err := tokens.GenerateAccessToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		refreshToken, err := tokens.GenerateRefreshToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		LogedUser.Adminemail = "admin@gmail.com"

		c.JSON(http.StatusOK, gin.H{"email": LogedUser.Adminemail, "access_token": accessToken,
			"refresh_token": refreshToken})

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}
