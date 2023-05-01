package route

import (
	"motor/controllers"
	"motor/middlewares"
	"motor/repositories/repoimpl"
	"motor/services/serviceimpl"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

func CreateHandler(db *gorm.DB) {
	r := gin.New()

	userRepo := repoimpl.CreateUserRepository(db)

	userService := serviceimpl.CreateUserService(userRepo)
	authService := serviceimpl.CreateAuthService(userRepo)

	r.Use(middlewares.SetupCorsMiddleware())

	api := r.Group("/api")
	controllers.CreateAuthRoutes(api.Group("/auth"), authService, userService)

	err := r.Run()

	if err != nil {
		log.Fatal(err)
		return
	}
}
