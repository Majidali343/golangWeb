package filemanipulate

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
		// Create a token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"Email": admin.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString([]byte(SecretKey))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error signing the token"})
			return
		}
		LogedUser.Adminemail = "admin@gmail.com"

		c.JSON(http.StatusOK, gin.H{"token": tokenString, "email": LogedUser.Adminemail})

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}
