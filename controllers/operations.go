package controllers

import (
	"net/http"
	"path/filepath"
	"encoding/base64"
	"strconv"
	"os"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	
	"../dbutils"
	"../types"
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

	maxDownloads, err := strconv.Atoi(c.PostForm("max_downloads"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	  	return
	}

	preferredUrl := c.PostForm("preferred_url")

	if preferredUrl == "" {
		preferredUrl = genUniqueName()
	} else {
		var fileinfos []types.DbModal
		dbutils.DB.Where("url = ?", preferredUrl).Find(&fileinfos)

		if len(fileinfos) > 0 {
			preferredUrl = genUniqueName()
		}
	}

	fileName := filepath.Base(preferredUrl + file.Filename)

	fmt.Printf("Preferred URL: %s Max Downloads: %d File_Name: %s\n", preferredUrl, maxDownloads, fileName)

	// Upload the file to specific dst.
	dst := save_dir + fileName
	c.SaveUploadedFile(file, dst)

	// save fileinfo
	fileInfo := types.DbModal{FileName: fileName, MaxDownloads: maxDownloads, Url: preferredUrl}
	dbutils.DB.Create(&fileInfo)

	c.JSON(http.StatusOK, gin.H{"url": preferredUrl})
}

func ServeFile(c *gin.Context) {
	fileStore := "./filestore/"

	url := c.Param("url")
	var fileInfo types.DbModal
	result := dbutils.DB.Where("url = ?", url).First(&fileInfo)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s", result.Error)})
	  	return
	}

	filepath := fileStore + fileInfo.FileName

	c.Header("Content-Description", "File Transfer")
    c.Header("Content-Transfer-Encoding", "binary")
    c.Header("Content-Disposition", "attachment; filename="+fileInfo.FileName)
	c.Header("Content-Type", "application/octet-stream")

	c.File(filepath)

	if fileInfo.MaxDownloads > 1 {
		fileInfo.MaxDownloads -= 1
		dbutils.DB.Save(&fileInfo)
	} else {
		dbutils.DB.Delete(&fileInfo)
		os.Remove(filepath)
	}
}

func genUniqueName() string {
	newUUID := uuid.Must(uuid.NewV4())
	uuidString := base64.RawURLEncoding.EncodeToString(newUUID.Bytes())
	return uuidString
}
