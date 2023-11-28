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
}

type MyUser struct {
	Email    string `json:"email" `
	Password string `json:"password"`
}

type logeduser struct {
	LoggedUserID uint   `json:"loggeduserid"`
	Adminemail   string `json:"email" `
}

var LogedUser logeduser

func Register(c *gin.Context) {
	var user MyUser

	// Parse JSON request into MyUser struct
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Connect to the database
	db := dbconnect.Dbconnection()
	defer db.Close()

	var existingUser User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		// User already exists
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}
	// Create a new user in the database
	var newUser User
	newUser.Email = user.Email
	newUser.Password = user.Password

	if err := db.Create(&newUser).Error; err != nil {
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

	var dbUser User

	// Connect to the database
	db := dbconnect.Dbconnection()
	defer db.Close()

	// Check if the user with the provided email exists
	if err := db.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify the password
	if dbUser.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

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

	LogedUser.LoggedUserID = dbUser.ID

	// Respond with user ID and token
	c.JSON(http.StatusOK, gin.H{"user_id": dbUser.ID, "token": tokenString})
}
