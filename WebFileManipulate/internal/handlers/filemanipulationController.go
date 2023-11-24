package filemanipulate

import (
	"fmt"
	"mime/multipart"
	"net/http"

	

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



var TotalCalculation calculation.Calculation
var ElapsedTime int64

func Filemanupulate(c *gin.Context) {
	var filedata filedata

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
}

func Details(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Calculation": TotalCalculation,
		"Time":        ElapsedTime,
	})
}
