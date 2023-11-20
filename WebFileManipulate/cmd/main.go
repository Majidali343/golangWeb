package main

import (
	"fmt"
	"net/http"

	"time"
	"wordcount/internal/calculation"
	"wordcount/internal/file"
	"wordcount/pkg/counting"

	"github.com/gin-gonic/gin"
)

type filedata struct {
	FileName string `form:"fileName" json:"fileName" binding:"required"`
	Routines int    `form:"routines" json:"routines" binding:"required"`
}

func filemanupulate() {
	router := gin.Default()

	var TotalCalculation calculation.Calculation
	var ElapsedTime int64

	router.POST("/filemanipulate", func(c *gin.Context) {
		var filedata filedata

		if err := c.ShouldBindJSON(&filedata); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		startTime := time.Now()

		data, err := file.ReadFile("../assets/" + filedata.FileName)
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

func main() {
	filemanupulate()
}
