package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type filedata struct {
	FileName string `form:"fileName" json:"fileName" binding:"required"`
	Routines int    `form:"routines" json:"routines" binding:"required"`
}

func filemanupulate() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, Gin!",
		})
	})

	router.POST("/file", func(c *gin.Context) {
		var requestBody filedata

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// startTime := time.Now()

		// data, err := file.ReadFile("../assets/" + filedata.FileName)
		// if err != nil {
		// 	fmt.Println("Error reading file:", err)
		// 	return
		// }

		// segmentSize := len(data) / filedata.Routines

		// doneCh := make(chan struct{})

		// partialResultCh := make(chan calculation.Calculation)

		// go func() {
		// 	var totalCalculation calculation.Calculation

		// 	for i := 0; i < requestBody.Routines; i++ {
		// 		partialResult := <-partialResultCh
		// 		totalCalculation.PunctuationCount += partialResult.PunctuationCount
		// 		totalCalculation.VowelCount += partialResult.VowelCount
		// 		totalCalculation.WordCount += partialResult.WordCount
		// 		totalCalculation.LineCount += partialResult.LineCount
		// 	}

		// 	fmt.Printf("Total details are %+v \n", totalCalculation)

		// 	close(doneCh)
		// }()

		// for i := 0; i < filedata.Routines; i++ {
		// 	go counting.Count(data[i*segmentSize:(i+1)*segmentSize], partialResultCh, doneCh)
		// }

		// <-doneCh

		// endTime := time.Now()
		// elapsedTime := endTime.Sub(startTime).Milliseconds()
		// fmt.Printf("Elapsed time: %d ms\n", elapsedTime)

		c.JSON(http.StatusOK, gin.H{"filedata": requestBody})
	})

	// Start the Gin server
	router.Run(":8080")
}

func main() {
	filemanupulate()
}
