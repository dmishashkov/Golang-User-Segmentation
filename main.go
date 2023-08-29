package main

import (
	"fmt"
	"github.com/dmishashkov/avito_test_task_2023/config"
	_ "github.com/dmishashkov/avito_test_task_2023/config"
	"github.com/dmishashkov/avito_test_task_2023/internal/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/getUserSlugs", controllers.GetSegments)
	r.POST("/createSlug", controllers.CreateSegment)
	r.PUT("/addRemoveSlugUser", controllers.UserSegments)
	r.DELETE("/deleteSlug", controllers.DeleteSegment)
	r.GET("/test", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"OK": "ok",
		})
	})
	r.Run(fmt.Sprintf("0.0.0.0:%d", config.ProjectConfig.Deploy.Port))
}
