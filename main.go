package main

import (
	"github.com/dmishashkov/avito_test_task_2023/internal/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/getUserSlugs", controllers.GetSegments)
	r.POST("/createSlug", controllers.CreateSegment)
	r.DELETE("/addRemoveSlugUser", controllers.UserSegments)
	r.DELETE("/deleteSlug", controllers.DeleteSegment)
	r.Run(":5050") // get from cfg
}
