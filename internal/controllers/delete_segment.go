package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dmishashkov/avito_test_task_2023/internal/db"
	"github.com/dmishashkov/avito_test_task_2023/internal/schemas"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteSegment(c *gin.Context) {
	database, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: "Error connecting to db",
		})
		return
	}
	request := struct {
		SegmentName string `json:"segment_name" binding:"required"`
	}{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: err.Error(),
		})
		return
	}
	statement := `SELECT segment_id FROM segments WHERE segment_name = $1`
	row := database.QueryRow(statement, request.SegmentName)
	val := 0
	if err := row.Scan(&val); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusConflict, schemas.Error{
				Error: "Nothing to delete: no such segment exists",
			})
		} else {
			c.JSON(http.StatusInternalServerError, schemas.Error{
				Error: fmt.Sprintf("%s: %s", "DB error", err.Error()),
			})
		}

		return
	}
	tr, err := database.Begin()
	defer tr.Rollback()
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: "Error starting transaction",
		})
		return
	}
	statement1 := `DELETE FROM segments WHERE segment_id = $1`
	statement2 := `DELETE FROM segments_users WHERE segment_id = $1`
	statement3 := `DELETE FROM segments_history WHERE segment_id = $1`

	_, err = tr.Exec(statement2, val)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: fmt.Sprintf("%s: %s", "DB error while deleting segment", err.Error()),
		})
		return
	}
	_, err = tr.Exec(statement1, val)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: fmt.Sprintf("%s: %s", "DB error while deleting segment", err.Error()),
		})
		return
	}
	_, err = tr.Exec(statement3, val)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: fmt.Sprintf("%s: %s", "DB error while deleting segment", err.Error()),
		})
		return
	}
	err = tr.Commit()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Error commiting transaction",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Operation executed successfully",
	})
}
