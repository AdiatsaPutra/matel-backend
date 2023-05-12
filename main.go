package main

import (
	"motor/controllers"
	"motor/security"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	r.GET("/profil", security.AuthMiddleware(), controllers.GetProfile)

	r.GET("/member", security.AuthMiddleware(), controllers.GetAllMember)

	r.GET("/province", controllers.GetProvince)
	r.GET("/kabupaten/:province-id", controllers.GetKabupaten)
	r.GET("/kecamatan/:kabupaten-id", controllers.GetKecamatan)

	r.Run()
}
