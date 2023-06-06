package controllers

import (
	"github.com/gin-gonic/gin"
)

func DownloadLeasing(c *gin.Context) {
	zipFilePath := "archive.zip"

	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+zipFilePath)
	c.Writer.Header().Set("Content-Type", "application/zip")

	c.File(zipFilePath)

	// Remove file after sending it to the user.
	// os.Remove(zipFilePath)

}
