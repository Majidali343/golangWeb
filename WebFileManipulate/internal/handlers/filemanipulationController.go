package filemanipulate

import (
	"fmt"

	"mime/multipart"
	"net/http"

	"path/filepath"
	"strconv"
	"time"
	"wordcount/internal/calculation"
	dbconnect "wordcount/internal/db"
	"wordcount/internal/filereader"
	"wordcount/pkg/counting"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type filedata struct {
	Routines int                   `form:"routines" json:"routines" binding:"required"`
	file     *multipart.FileHeader `form:"file"  binding:"required"`
}

type Filestatic struct {
	gorm.Model
	UserID   uint   `json:"user_id"`
	Filename string `json:"filename"`
	Time     int64  `json:"time"`
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

	db := dbconnect.Dbconnection()
	defer db.Close()
	db.AutoMigrate(&Filestatic{})

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

	Logeduserdetail := LogedUser.LoggedUserID

	var filestatic Filestatic
	filestatic.UserID = Logeduserdetail
	filestatic.Filename = filename
	filestatic.Time = ElapsedTime

	if err := db.Create(&filestatic).Error; err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"filedata": TotalCalculation, "Time": ElapsedTime, "logeduser": Logeduserdetail, "filename": filename})
}

var results []struct {
	Filename string  `json:"filename"`
	AvgTime  float64 `json:"avg_time"`
	Count    int64   `json:"count"`
}

func UserFileStatics(c *gin.Context) {

	db := dbconnect.Dbconnection()
	Logeduserdetail := LogedUser.LoggedUserID

	db.Table("filestatics").
		Select("filename, AVG(time) as avg_time, COUNT(*) as count").
		Where("user_id = ?", Logeduserdetail).
		Group("filename").
		Scan(&results)

	for _, result := range results {

		c.IndentedJSON(http.StatusOK, gin.H{
			"User ID":           Logeduserdetail,
			"Filename":          result.Filename,
			"File Average Time": result.AvgTime,
			"Counts":            result.Count,
		})

	}

}

func Admingetresults(c *gin.Context) {

	var results []struct {
		UserID   uint    `json:"user_id"`
		Filename string  `json:"filename"`
		AvgTime  float64 `json:"avg_time"`
		Count    int64   `json:"count"`
	}

	validadminaccess := LogedUser.Adminemail

	if validadminaccess != "" && validadminaccess == "admin@gmail.com" {

		db := dbconnect.Dbconnection()

		db.Table("filestatics").
			Select("user_id, filename, AVG(time) as avg_time, COUNT(*) as count").
			Group("user_id, filename").
			Scan(&results)

		for _, result := range results {

			c.IndentedJSON(http.StatusOK, gin.H{
				"Filename":          result.Filename,
				"File Average Time": result.AvgTime,
				"Counts":            result.Count,
				"User ID":           result.UserID,
			})

		}

	} else {

		c.JSON(http.StatusOK, gin.H{
			"Access error": "You can not access data you are a user ",
			"email":        validadminaccess,
		})
	}

}

func Details(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Calculation": TotalCalculation,
		"Time":        ElapsedTime,
	})
}
