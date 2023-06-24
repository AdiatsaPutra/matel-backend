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
	r.PATCH("/set-status", security.AuthMiddleware(), controllers.SetUser)

	r.GET("/profil", security.AuthMiddleware(), controllers.GetProfile)

	r.GET("/home", controllers.GetHome)

	r.GET("/leasing-master", controllers.GetLeasingMaster)
	r.POST("/leasing-master", controllers.CreateLeasing)
	r.PUT("/leasing-master/:id", controllers.UpdateLeasing)
	r.DELETE("/leasing-master/:id", controllers.DeleteLeasing)

	r.GET("/cabang", controllers.GetCabang)
	r.POST("/cabang", controllers.CreateCabang)
	r.PUT("/cabang/:id", controllers.UpdateCabang)
	r.DELETE("/cabang/:id", controllers.DeleteCabang)

	r.GET("/kendaraan", controllers.GetKendaraan)
	r.GET("/download-template", controllers.DownloadTemplate)
	r.DELETE("/delete-kendaraan", controllers.DeleteKendaraan)

	r.GET("/leasing", controllers.GetLeasing)
	r.GET("/leasing/:id", security.AuthMiddleware(), controllers.GetLeasingDetail)
	r.GET("/leasing/history", security.AuthMiddleware(), controllers.GetLeasingHistory)
	r.POST("/upload-leasing", controllers.AddCSV)
	r.GET("/dump-sql", controllers.DumpSQLHandler)
	r.POST("/download-update", controllers.UpdateSQLHandler)
	r.GET("/download-all", controllers.DownloadLeasing)
	r.GET("/download-app", controllers.DownloadApk)

	r.GET("/member", security.AuthMiddleware(), controllers.GetMember)
	r.POST("/update-member", controllers.SetUser)

	r.GET("/province", controllers.GetProvince)
	r.GET("/kabupaten/:province-id", controllers.GetKabupaten)
	r.GET("/kecamatan/:kabupaten-id", controllers.GetKecamatan)

	r.Run()
}
