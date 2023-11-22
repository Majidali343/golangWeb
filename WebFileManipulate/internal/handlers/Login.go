package filemanipulate

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const SecretKey = "Majid ali"

// User represents the structure of a user.
type User struct {
	Username string `json:"username"  form:"username" binding:"required"`
	Password string `json:"password"  form:"passsword" binding:"required"`
}

func Login(c *gin.Context) {

	var user User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if user.Username == "majid" && user.Password == "12345" {
		// Create a token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString([]byte(SecretKey))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error signing the token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}
