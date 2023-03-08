package routes
import (
	"github.com/gin-gonic/gin"
	controller "louis/go_projects/controllers"
)

func AuthRouts(incomingRoutes *gin.Engine){
	incomingRoutes.POST("users/signup",controller.SignUp())
	incomingRoutes.POST("users/login",controller.Login())
}