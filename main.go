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
	r.GET("/getUserSegments", controllers.GetSegments)
	r.GET("/getSegmentsHistory", controllers.GetHistory)
	r.POST("/createSegment", controllers.CreateSegment)
	r.PUT("/editUserSegments", controllers.UserSegments)
	r.DELETE("/deleteSegment", controllers.DeleteSegment)
	r.Run(fmt.Sprintf("0.0.0.0:%d", config.ProjectConfig.Deploy.Port))
}
