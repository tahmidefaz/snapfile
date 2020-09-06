package controllers

import (
	"net/http"
	"path/filepath"
	"github.com/gin-gonic/gin"
	"../dbutils"
	"../types"
	"strconv"
	"fmt"
)

func GetAll(c *gin.Context) {
	var fileinfos []types.DbModal
  	dbutils.DB.Find(&fileinfos)

	c.JSON(http.StatusOK, gin.H{
		"data": fileinfos,
	})
}


func UploadFilename(c *gin.Context) {
	// Validate input
	var input types.FileInfo
	if err := c.ShouldBindJSON(&input); err != nil {
	  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	  return
	}
  
	// Create a filename
	fileInfo := types.DbModal{FileName: input.FileName, MaxDownloads: input.MaxDownloads}
	dbutils.DB.Create(&fileInfo)
  
	c.JSON(http.StatusOK, gin.H{"data": input})
}

func UploadFile(c *gin.Context) {
	save_dir := "./filestore/"

	// Validate input
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	  	return
	}

	maxDownloads, _ := strconv.Atoi(c.PostForm("max_downloads"))
	preferredUrl := c.PostForm("preferred_url")

	if preferredUrl == "" {
		preferredUrl = genUniqueName()
	} else {
		var fileinfos []types.DbModal
		dbutils.DB.Where("file_name <> ?", preferredUrl).Find(&fileinfos)

		if len(fileinfos) > 0 {
			preferredUrl = genUniqueName()
		}

		fmt.Printf("Actual Url: %s", preferredUrl)
	}

	fileName := preferredUrl + filepath.Ext(file.Filename)

	fmt.Println(file.Filename)
	fmt.Printf("Preferred URL: %s Max Downloads: %d File_Name: %s\n", preferredUrl, maxDownloads, fileName)

	// Upload the file to specific dst.
	dst := save_dir + fileName
	c.SaveUploadedFile(file, dst)

	// save fileinfo
	fileInfo := types.DbModal{FileName: fileName, MaxDownloads: maxDownloads}
	dbutils.DB.Create(&fileInfo)

	message := fmt.Sprintf("%s uploaded successfully", file.Filename)
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func genUniqueName() string {
	return "unique_name"
}
