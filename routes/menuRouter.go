package routes

import (
	controllers "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/menus", controllers.GetMenus())
	incomingRoutes.GET("/menus/:menu_id", controllers.GetMenu())
	incomingRoutes.POST("/menus", controllers.CreateMenu())
	incomingRoutes.PATCH("/menus/:menu_id", controllers.UpdateMenu())
	incomingRoutes.DELETE("/menus/:menu_id", controllers.DeleteMenu())
}
