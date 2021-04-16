package routers

import (
	"github.com/gin-gonic/gin"
	"specGo2/examples/crudV7/controllers"
)

// SetupRouter ...
func SetupRouter() *gin.Engine {
	r := gin.Default()
	group1 := r.Group("/api")
	group1.GET("user", controllers.GetUsers)
	group1.GET("user/:id", controllers.GetUserById)
	group1.POST("user", controllers.CreateUser)
	group1.PUT("user/:id", controllers.UpdateUser)
	group1.DELETE("user/:id", controllers.DeleteUser)

	return r
}
