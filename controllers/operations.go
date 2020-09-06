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
	
	"github.com/tahmidefaz/snapfile/dbutils"
	"github.com/tahmidefaz/snapfile/types"
	"github.com/tahmidefaz/snapfile/misc"
)

var apiFilestore = misc.GetEnv("API_FILESTORE", "./filestore/")

func UploadFile(c *gin.Context) {
	save_dir := apiFilestore

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
	fileStore := apiFilestore

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

func PingHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func genUniqueName() string {
	newUUID := uuid.Must(uuid.NewV4())
	uuidString := base64.RawURLEncoding.EncodeToString(newUUID.Bytes())
	return uuidString
}
