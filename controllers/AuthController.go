package controllers

import (
	"errors"
	"motor/exceptions"
	"motor/payloads"
	"motor/security"
	"motor/services"
	"motor/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
}

func CreateAuthRoutes(r *gin.RouterGroup, authService services.AuthService, userService services.UserService) {
	authHandler := AuthController{
		authService: authService,
		userService: userService,
	}

	r.POST("/register", authHandler.DoRegister)
	r.POST("/login", authHandler.DoLogin)
}

func (r *AuthController) DoRegister(c *gin.Context) {
	var register payloads.CreateRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	check := utils.ValidationForm(register)

	if check != "" {
		exceptions.BadRequestException(c, check)
		return
	}

	findUser, _ := r.userService.FindUser(register.UserName)

	if findUser.UserName != "" {
		exceptions.NotFoundException(c, errors.New("Username already exists").Error())
		return
	}

	hash, err := security.HashPassword(register.Password)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	token, err := security.GenerateToken(findUser.UserName)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	register.Password = hash
	register.Token = token

	get, err := r.authService.DoRegister(register)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	payloads.HandleSuccess(c, get, "Register Successfully", http.StatusOK)
}

func (r *AuthController) DoLogin(c *gin.Context) {
	var login payloads.LoginRequest

	if err := c.ShouldBindJSON(&login); err != nil {
		exceptions.EntityException(c, err.Error())
		return
	}

	check := utils.ValidationForm(login)

	if check != "" {
		exceptions.BadRequestException(c, check)
		return
	}

	findUser, _ := r.userService.FindUser(login.Username)

	if findUser.UserName != "" {
		hashPwd := findUser.Password
		pwd := login.Password

		hash := security.VerifyPassword(hashPwd, pwd)

		if hash == nil {
			token, err := security.GenerateToken(findUser.UserName)

			if err != nil {
				exceptions.AppException(c, err.Error())
				return
			}

			findUser.Token = token

			payloads.HandleSuccess(c, findUser, "Login Successfully", http.StatusOK)
		} else {
			exceptions.BadRequestException(c, errors.New("Password dont matched").Error())
			return
		}
	} else {
		exceptions.NotFoundException(c, errors.New("Username not found").Error())
		return
	}
}
