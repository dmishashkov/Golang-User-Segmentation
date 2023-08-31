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

func GetSegments(c *gin.Context) {
	database, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: "Error connecting to db",
		})
		return
	}
	userID := struct {
		UserID *int `json:"user_id" binding:"required"`
	}{}
	err = c.BindJSON(&userID) // TODO: add HTTP status codes everywhere
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: err.Error(),
		})
		return
	}
	statement := `SELECT users.user_id FROM users WHERE
	                               user_id = $1`
	row := database.QueryRow(statement, userID.UserID)
	val := 0
	if err := row.Scan(&val); errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: "User with given ID does not exist",
		})
		return
	}
	statement = `SELECT segments.segment_name FROM segments
	JOIN segments_users USING(segment_id)
	WHERE user_id = $1`
	rows, err := database.Query(statement, userID.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: fmt.Sprintf("%s: %s", "DB error while inserting segment to random users", err.Error()),
		})
		return
	}
	segments := make([]string, 0)
	for rows.Next() {
		var t string
		rows.Scan(&t)
		segments = append(segments, t)
	}
	c.JSON(http.StatusOK, gin.H{
		"segments": segments,
	})
	defer rows.Close()
}
