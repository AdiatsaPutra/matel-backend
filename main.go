package main

import (
	"motor/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.Run()
}
