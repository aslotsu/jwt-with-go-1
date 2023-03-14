package routes

import (
	"github.com/gin-gonic/gin"
	controller "jwt-with-go/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.SignUp())
	incomingRoutes.POST("users/login", controller.Login())

}
