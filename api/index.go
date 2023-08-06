package api

import (
	"autos/controller"
	"autos/middleware"
	"autos/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	authController  controller.AuthController  = controller.NewAuthController(service.NewAuthService())
	autosController controller.AutosController = controller.NewAutosController(service.NewAutosService())
)

func Handler(w http.ResponseWriter, r *http.Request) {
	server := gin.Default()

	server.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/cars")
	})

	server.POST("/login", authController.Login)
	server.POST("/register", authController.Register)
	server.GET("/logout", middleware.TokenAuthMiddleware(true), authController.Logout)

	server.GET("/cars", middleware.TokenAuthMiddleware(false), autosController.ReadCars)
	server.POST("/create-car", middleware.TokenAuthMiddleware(true), autosController.CreateCar)
	server.PATCH("/edit-car", middleware.TokenAuthMiddleware(true), autosController.UpdateCar)
	server.DELETE("/delete-car", middleware.TokenAuthMiddleware(true), autosController.DeleteCar)

	server.ServeHTTP(w, r)
}
