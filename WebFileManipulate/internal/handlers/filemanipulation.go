package filemanipulate

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"path/filepath"
	"strconv"
	"time"
	"wordcount/internal/calculation"
	"wordcount/internal/filereader"
	"wordcount/pkg/counting"

	"github.com/gin-gonic/gin"
)

type filedata struct {
	Routines int                   `form:"routines" json:"routines" binding:"required"`
	file     *multipart.FileHeader `form:"file"  binding:"required"`
}

const SecretKey = "Majid ali"

// User represents the structure of a user.
type User struct {
	Username string `json:"username"  form:"username" binding:"required"`
	Password string `json:"password"  form:"passsword" binding:"required"`
}

func Filemanupulate() {
	router := gin.Default()

	router.POST("/login", func(c *gin.Context) {
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
	})

	var TotalCalculation calculation.Calculation
	var ElapsedTime int64

	router.POST("/filemanipulate", func(c *gin.Context) {
		var filedata filedata

		tokenString := c.GetHeader("Authorization")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret key
			return []byte(SecretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		routines := c.PostForm("routines")
		routinesInt, err := strconv.Atoi(routines)

		if err != nil {
			c.JSON(400, gin.H{"error": "No Routines  provided"})
			return
		}

		filedata.Routines = routinesInt

		///upload file feature ////////
		File, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "No file provided"})
			return
		}
		filedata.file = File

		filename := filepath.Base(filedata.file.Filename)

		// Save the uploaded file to the server
		err = c.SaveUploadedFile(filedata.file, "../assets/"+filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "File uploaded and read successfully",
			"filename": filename,
		})
		///end of upload file feature

		startTime := time.Now()

		data, err := filereader.ReadFile("../assets/" + filename)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		segmentSize := len(data) / filedata.Routines

		doneCh := make(chan struct{})

		partialResultCh := make(chan calculation.Calculation)

		go func() {

			for i := 0; i < filedata.Routines; i++ {
				partialResult := <-partialResultCh
				TotalCalculation.PunctuationCount += partialResult.PunctuationCount
				TotalCalculation.VowelCount += partialResult.VowelCount
				TotalCalculation.WordCount += partialResult.WordCount
				TotalCalculation.LineCount += partialResult.LineCount
			}

			// fmt.Printf("Total details are %+v \n", TotalCalculation)

			close(doneCh)
		}()

		for i := 0; i < filedata.Routines; i++ {
			go counting.Count(data[i*segmentSize:(i+1)*segmentSize], partialResultCh, doneCh)
		}

		<-doneCh

		endTime := time.Now()
		ElapsedTime = endTime.Sub(startTime).Milliseconds()

		c.JSON(http.StatusOK, gin.H{"filedata": TotalCalculation})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Calculation": TotalCalculation,
			"Time":        ElapsedTime,
		})
	})

	// Start the Gin server
	router.Run(":8080")
}
