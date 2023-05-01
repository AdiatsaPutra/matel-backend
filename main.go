package main

import (
	"motor/configs"
	"motor/controllers"

	"github.com/gin-gonic/gin"
)

func init() {
	configs.ConnectDatabase()
}

func main() {
	r := gin.Default()
	r.POST("/register", controllers.Register)
	r.Run()
}
