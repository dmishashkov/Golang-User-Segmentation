package controllers

import (
	"database/sql"
	"errors"
	"github.com/dmishashkov/avito_test_task_2023/internal/db"
	"github.com/dmishashkov/avito_test_task_2023/internal/schemas"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteSegment(c *gin.Context) {
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
	if err := row.Scan(&val); errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: "Nothing to delete: no such slug exists",
		})
		return
	}
	statement1 := `DELETE FROM slugs WHERE slug_id = $1`
	statement2 := `DELETE FROM slugs_users WHERE slug_id = $1;`

	_, err := database.Exec(statement2, val)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: err.Error(),
		})
		return
	}
	_, err = database.Exec(statement1, val)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"message": "Successfull deleted slug",
	})
}
