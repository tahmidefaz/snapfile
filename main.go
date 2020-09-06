package main

import (
	"github.com/gin-gonic/gin"
	"./dbutils"
	c "./controllers"
)

func main() {
	r := gin.Default()

	r.MaxMultipartMemory = 50 << 20  // 50 MiB

	dbutils.DatabaseConnect()

	v1 := r.Group("/v1")
	{
		v1.GET("/ping", c.PingHealth)
		v1.GET("/getall", c.GetAll)
		v1.POST("/", c.UploadFilename)
		v1.POST("/upload", c.UploadFile)
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
