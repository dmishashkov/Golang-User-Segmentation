package controllers

import (
	"fmt"
	"github.com/dmishashkov/avito_test_task_2023/internal/db"
	"github.com/dmishashkov/avito_test_task_2023/internal/schemas"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

func CreateSegment(c *gin.Context) {
	database := db.GetDB()
	segment := struct {
		SegmentName    string `json:"slug_name" binding:"required"`
		RandomPercents int    `json:"random_add"`
	}{}
	if err := c.BindJSON(&segment); err != nil {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: "Error processing JSON data",
		})
		return
	}
	statement := `SELECT slug_id FROM slugs WHERE slug_name = $1`
	row := database.QueryRow(statement, segment.SegmentName) // TODO: переделать на ошибку pg
	val := 0
	err := row.Scan(&val)
	if err == nil {
		c.JSON(http.StatusConflict, schemas.Error{
			Error: "Slug with given name already exists",
		})
		return
	}
	statement = `INSERT INTO slugs (slug_name) VALUES ($1) RETURNING slug_id;`
	slug_id_inserted := 0
	res := database.QueryRow(statement, segment.SegmentName)
	err = res.Scan(&slug_id_inserted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: fmt.Sprintf("%s: %s", "DB error", err.Error()),
		})
		return
	}
	if segment.RandomPercents != 0 {
		s1 := `SELECT DISTINCT user_id FROM slugs_users`
		s2 := `SELECT COUNT(DISTINCT user_id) FROM slugs_users`
		number_of_users := 0
		rows, err := database.Query(s1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, schemas.Error{
				Error: fmt.Sprintf("%s: %s", "DB error", err.Error()),
			})
			return
		}
		row := database.QueryRow(s2)
		if err := row.Scan(&number_of_users); err != nil {
			c.JSON(http.StatusInternalServerError, schemas.Error{
				Error: fmt.Sprintf("%s: %s", "DB error", err.Error()),
			})
			return
		}

		if number_of_users == 0 {
			c.JSON(http.StatusConflict, schemas.Error{
				Error: "Can not add slugs because zero users exist",
			})
			return
		}
		total := int(math.Ceil(float64(number_of_users) / float64(100) * float64(segment.RandomPercents)))
		defer rows.Close()
		for rows.Next() {
			if total > 0 {
				user_id := 0
				rows.Scan(&user_id)
				s := `INSERT INTO slugs_users VALUES ($1, $2)`
				_, err := database.Exec(s, slug_id_inserted, user_id)
				if err != nil {
					c.JSON(http.StatusInternalServerError, schemas.Error{
						Error: fmt.Sprintf("%s: %s", "DB error while inserting slug to random users", err.Error()),
					})
					return
				}
				total--
			} else {
				break
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully created segment and inserted to given percents of random users",
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully created segment",
	})

}
