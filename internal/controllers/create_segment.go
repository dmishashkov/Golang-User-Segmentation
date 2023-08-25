package controllers

import (
	"github.com/dmishashkov/avito_test_task_2023/internal/db"
	"github.com/dmishashkov/avito_test_task_2023/internal/schemas"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateSegment(c *gin.Context) {
	database := db.GetDB()
	segment := struct {
		SegmentName string `json:"segment_name" binding:"required"`
	}{}
	if err := c.BindJSON(&segment); err != nil {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: err.Error(),
		})
		return
	}
	statement := `SELECT slug_id FROM slugs WHERE slug_name = $1`
	row := database.QueryRow(statement, segment.SegmentName)
	val := 0
	err := row.Scan(&val)
	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: "Slug with given name already exists",
		})
		return
	}
	statement = `INSERT INTO slugs VALUES ($1);`
	_, err = database.Exec(statement, segment.SegmentName)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully created statement",
	})

}
