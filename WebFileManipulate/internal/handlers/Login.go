package filemanipulate

import (
	"net/http"
	"time"
	dbconnect "wordcount/internal/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const SecretKey = "Majid ali"

// User represents the structure of a user.
type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"not null;unique"`
	Password string `json:"password" gorm:"not null"`
	// Add other fields as needed
}

type MyUser struct {
	Email    string `json:"email" `
	Password string `json:"password"`
	// Add other fields as needed
}

func Register(c *gin.Context) {

	var user MyUser
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbuser User
	dbuser.Email = user.Email
	dbuser.Password = user.Password

	db := dbconnect.Dbconnection()
	db.AutoMigrate(&User{})
	defer db.Close()
	// Create a new user in the database
	if err := db.Create(&dbuser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {

	var user MyUser

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if user.Email == "majid" && user.Password == "12345" {
		// Create a token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"Email": user.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
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
