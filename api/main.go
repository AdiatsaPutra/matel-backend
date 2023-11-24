package main

import (
	"matel/controllers"
	"matel/controllers/kendaraan"
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
	r.GET("/kendaraan-per-cabang", controllers.GetKendaraanPerCabang)
	r.GET("/kendaraan-total", controllers.GetTotalKendaraan)

	r.GET("/leasing-master", controllers.GetLeasingMaster)
	r.POST("/leasing-master", controllers.CreateLeasing)
	r.PUT("/leasing-master/:id", controllers.UpdateLeasing)
	r.DELETE("/leasing-master/:id", controllers.DeleteLeasing)

	r.GET("/cabang", controllers.GetCabang)
	r.GET("/cabang-with-total", controllers.GetCabangWithTotal)
	r.GET("/cabang-export", controllers.GetCabangExport)
	r.POST("/cabang", controllers.CreateCabang)
	r.PUT("/cabang/:id", controllers.UpdateCabang)
	r.DELETE("/cabang/:id", controllers.DeleteCabang)

	r.GET("/kendaraan", controllers.GetKendaraan)
	r.GET("/download-template", controllers.DownloadTemplate)
	r.GET("/download-template-cabang", controllers.DownloadTemplateCabang)
	r.DELETE("/delete-kendaraan", controllers.DeleteKendaraan)
	r.DELETE("/delete-kendaraan/:id", controllers.DeleteKendaraanByID)
	r.DELETE("/delete-all-kendaraan", controllers.DeleteAllKendaraan)

	r.GET("/leasing", controllers.GetLeasing)
	r.GET("/leasing/:id", security.AuthMiddleware(), controllers.GetLeasingDetail)
	r.GET("/leasing/history", security.AuthMiddleware(), controllers.GetLeasingHistory)
	r.POST("/upload-leasing", controllers.AddCSV)
	r.POST("/upload-leasing-per-cabang", kendaraan.AddCSVPerCabang)
	r.GET("/dump-sql", controllers.DumpSQLHandler)
	r.POST("/download-update", controllers.UpdateSQLHandler)
	r.GET("/download-all", controllers.DownloadLeasing)
	r.GET("/download-app", controllers.DownloadApk)

	r.GET("/member", security.AuthMiddleware(), controllers.GetMember)
	r.POST("/update-member", controllers.SetUser)
	r.DELETE("/delete-member/:id", controllers.DeleteMember)

	r.GET("/province", controllers.GetProvince)
	r.GET("/kabupaten/:province-id", controllers.GetKabupaten)
	r.GET("/kecamatan/:kabupaten-id", controllers.GetKecamatan)

	r.GET("/banks", controllers.GetBanks)
	r.GET("/bank-data", controllers.GetBankData)
	r.GET("/banks/:id", controllers.GetBank)
	r.POST("/banks", controllers.CreateBank)
	r.PUT("/banks/:id", controllers.UpdateBank)
	r.DELETE("/banks/:id", controllers.DeleteBank)

	r.POST("/version", controllers.GetLatestVersion)

	r.Run()
}
