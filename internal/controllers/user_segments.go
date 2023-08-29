package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dmishashkov/avito_test_task_2023/internal/db"
	"github.com/dmishashkov/avito_test_task_2023/internal/schemas"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"time"
)

func createDeleteFunc(userID int, slugID int) func() {
	return func() {
		database := db.GetDB()
		s := `DELETE FROM slugs_users WHERE user_id = $1 AND slug_id = $2`
		_, err := database.Exec(s, userID, slugID)
		fmt.Println(err)
	}
}
func UserSegments(c *gin.Context) {
	database := db.GetDB()
	request := struct {
		UserID      *int              `json:"user_id" binding:"required"`
		DeleteSlugs []string          `json:"delete_slugs"`
		AddSlugs    []schemas.AddSlug `json:"add_slugs"`
	}{
		DeleteSlugs: make([]string, 0),
		AddSlugs:    make([]schemas.AddSlug, 0),
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: "Error processing JSON data",
		})
		return
	}
	response := struct {
		DeletedSlugs []string
		AddedSlugs   []string
		Errors       []schemas.Error
	}{
		make([]string, 0), make([]string, 0), make([]schemas.Error, 0),
	}
	for _, el := range request.AddSlugs {
		s := `SELECT slug_id FROM slugs WHERE slug_name = $1`
		row := database.QueryRow(s, el.Name)
		slug_id := 0
		if err := row.Scan(&slug_id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				response.Errors = append(response.Errors, schemas.Error{fmt.Sprintf("Slug with name %s does not exist", el.Name)})
			} else {
				response.Errors = append(response.Errors, schemas.Error{err.Error()})
			}

		} else {
			s := `INSERT INTO slugs_users (user_id, slug_id) VALUES ($1, $2)`
			_, err = database.Exec(s, request.UserID, slug_id)
			if err != nil {
				if err.(*pq.Error).Code == "23505" {
					response.Errors = append(response.Errors, schemas.Error{fmt.Sprintf("User already in slug %s", el.Name)})
				} else {
					response.Errors = append(response.Errors, schemas.Error{err.Error()})
				}

			} else {
				response.AddedSlugs = append(response.AddedSlugs, el.Name)
			}
			if el.Time != nil {
				time.AfterFunc((*el.Time).Sub(time.Now()), createDeleteFunc(*request.UserID, slug_id))
			}
		}

	}
	for _, el := range request.DeleteSlugs {
		s := `SELECT slug_id FROM slugs WHERE slug_name = $1`
		row := database.QueryRow(s, el)
		slug_id := 0
		if err := row.Scan(&slug_id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				response.Errors = append(response.Errors, schemas.Error{fmt.Sprintf("Slug with name %s does not exist", el)})
			} else {
				response.Errors = append(response.Errors, schemas.Error{err.Error()})
			}
		} else {
			s := `DELETE FROM slugs_users WHERE slug_id = $1 AND user_id =  $2` // (slug_id, user_id)
			res, err := database.Exec(s, slug_id, request.UserID)
			if err != nil {
				response.Errors = append(response.Errors, schemas.Error{err.Error()})
			} else if affected, _ := res.RowsAffected(); affected == 0 {
				response.Errors = append(response.Errors, schemas.Error{fmt.Sprintf("User wasn't in slug with name %s", el)})
			} else {
				response.DeletedSlugs = append(response.DeletedSlugs, el)
			}
		}
	}
	c.JSON(http.StatusOK, response)

}
