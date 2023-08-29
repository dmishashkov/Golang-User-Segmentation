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
	database := db.GetDB()
	segment := struct {
		SegmentName string `json:"slug_name" binding:"required"`
	}{}
	if err := c.BindJSON(&segment); err != nil {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: "Error processing JSON data",
		})
		return
	}
	statement := `SELECT slug_id FROM slugs WHERE slug_name = $1`
	row := database.QueryRow(statement, segment.SegmentName)
	val := 0
	if err := row.Scan(&val); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusConflict, schemas.Error{
				Error: "Nothing to delete: no such slug exists",
			})
		} else {
			c.JSON(http.StatusInternalServerError, schemas.Error{
				Error: fmt.Sprintf("%s: %s", "DB error", err.Error()),
			})
		}

		return
	}
	statement1 := `DELETE FROM slugs WHERE slug_id = $1`
	statement2 := `DELETE FROM slugs_users WHERE slug_id = $1;`

	_, err := database.Exec(statement2, val)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: fmt.Sprintf("%s: %s", "DB error", err.Error()),
		})
		return
	}
	_, err = database.Exec(statement1, val)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: fmt.Sprintf("%s: %s", "DB error", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted slug",
	})
}
