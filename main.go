package main

import (
	"motor/controllers"
	"motor/middlewares"
	"motor/security"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(middlewares.SetupCorsMiddleware())

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	r.GET("/profil", security.AuthMiddleware(), controllers.GetProfile)

	r.GET("/leasing", controllers.GetLeasing)
	r.POST("/upload-leasing", controllers.AddCSV)
	r.GET("/export", controllers.ExportHandler)

	r.GET("/province", controllers.GetProvince)
	r.GET("/kabupaten/:province-id", controllers.GetKabupaten)
	r.GET("/kecamatan/:kabupaten-id", controllers.GetKecamatan)

	r.Run()
}
