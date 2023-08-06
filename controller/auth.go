package controller

import (
	"autos/entity"
	"autos/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type authController struct {
	service service.AuthService
}

func NewAuthController(service service.AuthService) AuthController {
	return &authController{
		service: service,
	}
}

func (controller *authController) Login(ctx *gin.Context) {
	var request entity.AuthRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Email and password are required",
		})
		return
	}

	token, err := controller.service.Login(request.Email, request.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	ctx.Header("Authorization", "Bearer "+token)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully logged in",
	})
}

func (controller *authController) Register(ctx *gin.Context) {
	var request entity.AuthRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Email and password are required",
		})
		return
	}

	err = controller.service.Register(request.Email, request.Password)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Successfully registered %s", request.Email),
	})
}

func (controller *authController) Logout(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	controller.service.Logout(token.(string))

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully logged out",
	})
}
