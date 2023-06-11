package main

import (
	"matel/controllers"
	"matel/middlewares"
	"matel/security"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(middlewares.SetupCorsMiddleware())

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/login-web", controllers.LoginWeb)
	r.PATCH("/logout", security.AuthMiddleware(), controllers.Logout)
	r.PATCH("/reset-device-id", security.AuthMiddleware(), controllers.ResetDeviceID)

	r.GET("/profil", security.AuthMiddleware(), controllers.GetProfile)

	r.GET("/home", controllers.GetHome)

	r.GET("/leasing", controllers.GetLeasing)
	r.GET("/leasing/:id", security.AuthMiddleware(), controllers.GetLeasingDetail)
	r.GET("/leasing/history", security.AuthMiddleware(), controllers.GetLeasingHistory)
	r.POST("/upload-leasing", controllers.AddCSV)
	r.GET("/dump-sql", controllers.DumpSQLHandler)
	r.GET("/download-update", controllers.UpdateSQLHandler)
	r.GET("/download-all", controllers.DownloadLeasing)

	r.GET("/member", security.AuthMiddleware(), controllers.GetMember)

	r.GET("/province", controllers.GetProvince)
	r.GET("/kabupaten/:province-id", controllers.GetKabupaten)
	r.GET("/kecamatan/:kabupaten-id", controllers.GetKecamatan)

	r.Run()
}
