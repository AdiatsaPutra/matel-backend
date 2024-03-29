package controllers

import (
	"path"

	"github.com/gin-gonic/gin"
)

func DownloadLeasing(c *gin.Context) {
	zipFilePath := "dump_all.zip"

	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+zipFilePath)
	c.Writer.Header().Set("Content-Type", "application/zip")

	c.File(zipFilePath)

	// Remove file after sending it to the user.
	// os.Remove(zipFilePath)

}

func DownloadApk(c *gin.Context) {
	dir := "app"
	fileApk := "app-release.apk"

	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fileApk)
	c.Writer.Header().Set("Content-Type", "application/vnd.android.package-archive")

	c.File(path.Join(dir, fileApk))

}
