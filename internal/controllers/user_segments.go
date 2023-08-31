package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dmishashkov/avito_test_task_2023/internal/db"
	"github.com/dmishashkov/avito_test_task_2023/internal/schemas"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

func createDeleteFunc(userID int, segmentID int) func() {
	return func() {
		database, err := db.GetDB()
		s := `DELETE FROM segments_users WHERE user_id = $1 AND segment_id = $2`
		_, err = database.Exec(s, userID, segmentID)
		if err != nil {
			log.Print(err.Error())
		}
		s1 := `INSERT INTO segments_history (user_id, segment_id, action_date, action_type) VALUES ($1, $2, $3, $4)`
		_, err = database.Exec(s1, userID, segmentID, time.Now().String(), "DELETED")
		if err != nil {
			log.Print(err.Error())
		}
	}
}
func UserSegments(c *gin.Context) {
	database, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: "Error connecting to db",
		})
		return
	}
	request := struct {
		UserID         *int                 `json:"user_id" binding:"required"`
		DeleteSegments []string             `json:"delete_segments"`
		AddSegments    []schemas.AddSegment `json:"add_segments"`
	}{
		DeleteSegments: make([]string, 0),
		AddSegments:    make([]schemas.AddSegment, 0),
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: err.Error(),
		})
		return
	}
	response := struct {
		DeletedSegments []string `json:"deleted_segments"`
		AddedSegments   []string `json:"added_segments"`
		Errors          []schemas.Error
	}{
		make([]string, 0), make([]string, 0), make([]schemas.Error, 0),
	}
	tr, err := database.Begin()
	defer tr.Rollback()
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: "Error starting transaction",
		})
		return
	}
	for _, el := range request.AddSegments {
		s := `SELECT segment_id FROM segments WHERE segment_name = $1`
		row := tr.QueryRow(s, el.Name)
		segmentId := 0
		if err := row.Scan(&segmentId); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				response.Errors = append(response.Errors, schemas.Error{fmt.Sprintf("segment with name %s does not exist", el.Name)})
			} else {
				c.JSON(http.StatusInternalServerError, schemas.Error{
					Error: fmt.Sprintf("%s: %s", "DB error", err.Error()),
				})
				return
			}

		} else {
			s := `INSERT INTO segments_users (user_id, segment_id) VALUES ($1, $2)`
			_, err = tr.Exec(s, request.UserID, segmentId)
			if err != nil {
				if err.(*pq.Error).Code == "23505" {
					response.Errors = append(response.Errors, schemas.Error{fmt.Sprintf("User already in segment %s", el.Name)})
				} else {
					c.JSON(http.StatusInternalServerError, schemas.Error{
						Error: fmt.Sprintf("%s: %s", "DB error", err.Error()),
					})
					return
				}

			} else {
				response.AddedSegments = append(response.AddedSegments, el.Name)
				s1 := `INSERT INTO segments_history (user_id, segment_id, action_date, action_type) VALUES ($1, $2, $3, $4)`
				_, err := tr.Exec(s1, request.UserID, segmentId, time.Now(), "ADDED")
				if err != nil {
					c.JSON(http.StatusInternalServerError, schemas.Error{
						Error: fmt.Sprintf("%s: %s", "DB error", err.Error()),
					})
					return
				}
			}

			if el.Time != nil {
				time.AfterFunc((*el.Time).Local().Sub(time.Now()), createDeleteFunc(*request.UserID, segmentId))
			}
		}

	}
	for _, el := range request.DeleteSegments {
		s := `SELECT segment_id FROM segments WHERE segment_name = $1`
		row := tr.QueryRow(s, el)
		segmentId := 0
		if err := row.Scan(&segmentId); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				response.Errors = append(response.Errors, schemas.Error{fmt.Sprintf("segment with name %s does not exist", el)})
			} else {
				c.JSON(http.StatusInternalServerError, schemas.Error{
					Error: fmt.Sprintf("%s: %s", "DB error", err.Error()),
				})
				return
			}
		} else {
			s := `DELETE FROM segments_users WHERE segment_id = $1 AND user_id =  $2` // (segment_id, user_id)
			res, err := tr.Exec(s, segmentId, request.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, schemas.Error{
					Error: fmt.Sprintf("%s: %s", "DB error", err.Error()),
				})
				return
			} else if affected, _ := res.RowsAffected(); affected == 0 {
				response.Errors = append(response.Errors, schemas.Error{fmt.Sprintf("User wasn't in segment with name %s", el)})
			} else {
				response.DeletedSegments = append(response.DeletedSegments, el)
				s1 := `INSERT INTO segments_history (user_id, segment_id, action_date, action_type) VALUES ($1, $2, $3, $4)`
				_, err := tr.Exec(s1, request.UserID, segmentId, time.Now(), "DELETED")
				if err != nil {
					c.JSON(http.StatusInternalServerError, schemas.Error{
						Error: fmt.Sprintf("%s: %s", "DB error", err.Error()),
					})
					return
				}
			}
		}
	}
	err = tr.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: "Error committing transaction",
		})
		return
	}
	c.JSON(http.StatusOK, response)

}
