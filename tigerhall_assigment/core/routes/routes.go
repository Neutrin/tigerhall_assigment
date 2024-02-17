package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nitin/tigerhall/core/controllers"
	middlewares "github.com/nitin/tigerhall/core/middleware"
)

func Routes(r *gin.Engine, authController controllers.AuthController,
	tigerController controllers.TigerControllers) {
	//GET REQUEST
	r.GET("/tigers", tigerController.ListAllTigers)
	r.GET("/tigers/:tiger_id/sightings", tigerController.ListAllSightings)

	//POST REQUEST
	r.POST("/login", authController.Login)
	r.POST("/signup", authController.Signup)
	r.POST("/sighting", middlewares.IsAuthorized(), tigerController.AddTigerSighting)
	r.POST("/tigers", middlewares.IsAuthorized(), tigerController.AddTiger)

	//r.POST("/testing", middlewares.IsAuthorized(), testingController.ShowMessage)
	//r.POST("/testing", testingController.ShowMessage)
	//r.GET("/home", controllers.Home)
	//r.GET("/premium", controllers.Premium)
	//r.GET("/logout", controllers.Logout)
}
