package controller

import (
	"autos/autoserror"
	"autos/entity"
	"autos/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AutosController interface {
	ReadCars(ctx *gin.Context)
	CreateCar(ctx *gin.Context)
	UpdateCar(ctx *gin.Context)
	DeleteCar(ctx *gin.Context)
}

type autosController struct {
	service service.AutosService
}

func NewAutosController(service service.AutosService) AutosController {
	return &autosController{
		service: service,
	}
}

func (controller *autosController) ReadCars(ctx *gin.Context) {
	_, userId := getUserId(ctx)
	cars, err := controller.service.ReadCars(userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  http.StatusBadGateway,
			"message": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Fetched all cars",
		"data":    cars,
	})
}

func (controller *autosController) CreateCar(ctx *gin.Context) {
	_, userId := getUserId(ctx)

	var request entity.NewCarRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Name (string), company (string) and model (int) are required",
		})
		return
	}

	carId, err := controller.service.CreateCar(request, userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  http.StatusBadGateway,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Created new car " + carId,
	})
}

func (controller *autosController) UpdateCar(ctx *gin.Context) {
	_, userId := getUserId(ctx)

	var request entity.ModifyCarRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Car id (string) is required",
		})
		return
	}

	carId, err := controller.service.UpdateCar(request, userId)
	if customErr, ok := err.(*autoserror.CustomError); ok {
		ctx.JSON(customErr.Status, gin.H{
			"status":  customErr.Status,
			"message": customErr.Message,
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  http.StatusBadGateway,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Updated car " + carId,
	})
}

func (controller *autosController) DeleteCar(ctx *gin.Context) {
	_, userId := getUserId(ctx)

	var request entity.ModifyCarRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Car id (string) is required",
		})
		return
	}

	err = controller.service.DeleteCar(request.Id, userId)
	if customErr, ok := err.(*autoserror.CustomError); ok {
		ctx.JSON(customErr.Status, gin.H{
			"status":  customErr.Status,
			"message": customErr.Message,
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  http.StatusBadGateway,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Deleted car " + request.Id,
	})

}

func getUserId(ctx *gin.Context) (isAuthenticated bool, userId string) {
	authStatus, _ := ctx.Get("isAuthenticated")
	if authStatus.(bool) {
		id, _ := ctx.Get("userId")
		userId = id.(string)
	}
	return authStatus.(bool), userId
}
