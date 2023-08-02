package main

import (
	"autos/db"
	"autos/entity"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, "Hello, World")
	})

	server.POST("/add-car", postFunc)

	server.Run(":8080")
}

func postFunc(c *gin.Context) {
	var car entity.Car

	if err := c.BindJSON(&car); err != nil {
		fmt.Println("Error", err)
		return
	}
	fmt.Printf("Car is %T", car)

	id, err := db.InsertAuto(car)
	if err != nil {
		c.JSON(404, "Something went wrong")
	}

	c.JSON(200, id)

}
