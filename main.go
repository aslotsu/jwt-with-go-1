package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jwt-with-go/routes"

	//"jwt-with-go/routes"
	"log"
	"os"
)

func checkRouteOne(c *gin.Context) {
	c.JSON(200, gin.H{"Success": "Reached home route, no problem"})
}

func main() {

	port := os.Getenv("PORT")
	router := gin.New()
	router.Use(gin.Logger())
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/", checkRouteOne)
	fmt.Println("Crazy jwt auth api is working baby!")

	if err := router.Run(port); err != nil {
		log.Fatal(err)
	}

}
