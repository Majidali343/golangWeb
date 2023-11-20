package filemanipulate

import (
	"fmt"
	"net/http"

	"path/filepath"
	"strconv"
	"time"
	"wordcount/internal/calculation"
	"wordcount/internal/filereader"
	"wordcount/pkg/counting"

	"github.com/gin-gonic/gin"
)

// type filedata struct {
// 	Routines int `form:"routines" json:"routines" binding:"required"`
// 	// file     *multipart.FileHeader `form:"file" json:"file" binding:"required"`
// }

func Filemanupulate() {
	router := gin.Default()

	var TotalCalculation calculation.Calculation
	var ElapsedTime int64

	router.POST("/filemanipulate", func(c *gin.Context) {
		// var filedata filedata

		// if err := c.ShouldBindJSON(&filedata); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }

		// var Routines int = c.PostForm("routines")
		// Convert string to int for Routines

		routines := c.PostForm("routines")
		Routines, err := strconv.Atoi(routines)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid value for routines"})
			return
		}

		///upload file feature ////////

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filename := filepath.Base(file.Filename)
		// ext := filepath.Ext(filename)

		// Save the uploaded file to the server
		err = c.SaveUploadedFile(file, "../assets/"+filename)
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

		segmentSize := len(data) / Routines

		doneCh := make(chan struct{})

		partialResultCh := make(chan calculation.Calculation)

		go func() {

			for i := 0; i < Routines; i++ {
				partialResult := <-partialResultCh
				TotalCalculation.PunctuationCount += partialResult.PunctuationCount
				TotalCalculation.VowelCount += partialResult.VowelCount
				TotalCalculation.WordCount += partialResult.WordCount
				TotalCalculation.LineCount += partialResult.LineCount
			}

			// fmt.Printf("Total details are %+v \n", TotalCalculation)

			close(doneCh)
		}()

		for i := 0; i < Routines; i++ {
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
