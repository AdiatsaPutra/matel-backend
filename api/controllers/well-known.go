package controllers

import (
	"path"

	"github.com/gin-gonic/gin"
)

func AssetLinks(c *gin.Context) {
	dir := ".well-known"
	file := "assetlinks.json"

	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+file)
	c.Writer.Header().Set("Content-Type", "application/json")

	c.File(path.Join(dir, file))

}
