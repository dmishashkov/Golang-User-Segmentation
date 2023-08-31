package controllers

import (
	"fmt"
	"github.com/dmishashkov/avito_test_task_2023/internal/db"
	"github.com/dmishashkov/avito_test_task_2023/internal/schemas"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"time"
)

func CreateSegment(c *gin.Context) {
	database, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: "Error connecting to db",
		})
		return
	}
	request := struct {
		SegmentName    string `json:"segment_name" binding:"required"`
		RandomPercents int    `json:"random_add"`
	}{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: err.Error(),
		})
		return
	}
	if request.RandomPercents < 0 || request.RandomPercents > 100 {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: "Percent should be in range (0;100]",
		})
		return
	}
	s := `SELECT segment_id FROM segments WHERE segment_name = $1`
	row := database.QueryRow(s, request.SegmentName) // TODO: переделать на ошибку pg
	segmentId := 0
	err = row.Scan(&segmentId)
	if err == nil {
		c.JSON(http.StatusConflict, schemas.Error{
			Error: "Segment with given name already exists",
		})
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
	statement := `INSERT INTO segments (segment_name) VALUES ($1) RETURNING segment_id;`
	segmentIdInserted := 0
	res := tr.QueryRow(statement, request.SegmentName)
	err = res.Scan(&segmentIdInserted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: fmt.Sprintf("%s: %s", "DB error while creating segments", err.Error()),
		})
		return
	}
	if request.RandomPercents != 0 {
		s1 := `SELECT user_id FROM users`
		s2 := `SELECT COUNT(user_id) FROM users`
		numberOfUsers := 0
		row := tr.QueryRow(s2)
		if err := row.Scan(&numberOfUsers); err != nil {
			c.JSON(http.StatusInternalServerError, schemas.Error{
				Error: fmt.Sprintf("%s: %s", "DB error while creating segment to random users", err.Error()),
			})
			return
		}
		rows, err := tr.Query(s1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, schemas.Error{
				Error: fmt.Sprintf("%s: %s", "DB error while creating segment to random users", err.Error()),
			})
			return
		}

		if numberOfUsers == 0 {
			c.JSON(http.StatusConflict, schemas.Error{
				Error: "Can not add segments because zero users exist",
			})
			return
		}
		total := int(math.Ceil(float64(numberOfUsers) / float64(100) * float64(request.RandomPercents)))
		users := make([]int, 0)
		for rows.Next() {
			userID := 0
			rows.Scan(&userID)
			users = append(users, userID)
		}
		rows.Close()
		for _, userID := range users {
			if total > 0 {
				s := `INSERT INTO segments_users (segment_id, user_id) VALUES ($1, $2)`
				_, err := tr.Exec(s, segmentIdInserted, userID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, schemas.Error{
						Error: fmt.Sprintf("%s: %s", "DB error while inserting segment to random users", err.Error()),
					})
					return
				}
				s1 := `INSERT INTO segments_history (user_id, segment_id, action_date, action_type) VALUES ($1, $2, $3, $4)`
				_, err = tr.Exec(s1, userID, segmentIdInserted, time.Now(), "ADDED")
				if err != nil {
					c.JSON(http.StatusInternalServerError, schemas.Error{
						Error: fmt.Sprintf("%s: %s", "DB error while inserting segment to random users", err.Error()),
					})
					return
				}
				total--
			} else {
				break
			}
		}

	}
	err = tr.Commit()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Error committing transaction",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Operation executed successfully",
	})

}
