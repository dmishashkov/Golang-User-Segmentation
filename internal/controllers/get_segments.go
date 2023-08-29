package controllers

import (
	"fmt"
	"github.com/dmishashkov/avito_test_task_2023/internal/db"
	"github.com/dmishashkov/avito_test_task_2023/internal/schemas"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSegments(c *gin.Context) {
	database := db.GetDB()
	userID := struct {
		UserID *int `json:"user_id" binding:"required"`
	}{}
	err := c.BindJSON(&userID) // TODO: add HTTP status codes everywhere
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: "Error processing JSON data",
		})
		return
	}
	//statement := `SELECT users.user_id FROM users WHERE
	//                                user_id = $1`
	//row := database.QueryRow(statement, userID.UserID)
	//val := 0
	//if err := row.Scan(&val); errors.Is(err, sql.ErrNoRows) {
	//	c.JSON(http.StatusUnprocessableEntity, schemas.Error{
	//		Error: "User with given ID does not exist",
	//	})
	//	s := `INSERT INTO users VALUES ($1)`
	//	_, err = database.Exec(s, userID.UserID)
	//	return
	//}
	statement := `SELECT slugs.slug_name FROM slugs
	JOIN slugs_users USING(slug_id)
	WHERE user_id = $1`
	rows, err := database.Query(statement, userID.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: fmt.Sprintf("%s: %s", "DB error while inserting slug to random users", err.Error()),
		})
		return
	}
	slugs := make([]string, 0)
	for rows.Next() {
		var t string
		rows.Scan(&t)
		slugs = append(slugs, t)
	}
	c.JSON(http.StatusOK, gin.H{
		"slugs": slugs,
	})
	defer rows.Close()
}
