package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nitin/tigerhall/core/controllers"
)

func Routes(r *gin.Engine, authController controllers.AuthController,
	tigerController controllers.TigerControllers) {
	r.POST("/login", authController.Login)
	r.POST("/signup", authController.Signup)
	r.POST("/sighting", tigerController.AddTigerSighting)
	r.POST("/tigers", tigerController.AddTiger)
	r.GET("/tigers", tigerController.ListAllTigers)
	//r.POST("/testing", middlewares.IsAuthorized(), testingController.ShowMessage)
	//r.POST("/testing", testingController.ShowMessage)
	//r.GET("/home", controllers.Home)
	//r.GET("/premium", controllers.Premium)
	//r.GET("/logout", controllers.Logout)
}
