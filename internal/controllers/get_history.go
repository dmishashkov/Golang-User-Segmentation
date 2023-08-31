package controllers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/dmishashkov/avito_test_task_2023/internal/db"
	"github.com/dmishashkov/avito_test_task_2023/internal/schemas"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetHistory(c *gin.Context) {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	w.Comma = ';'
	database, err := db.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: "Error connecting to db",
		})
		return
	}
	//f, err := os.Create("../")
	//if err != nil {
	//	log.Fatal(err)
	//}
	request := struct {
		Begin *time.Time `json:"begin" binding:"required"`
		End   *time.Time `json:"end" binding:"required"`
	}{}
	err = c.BindJSON(&request) // TODO: add HTTP status codes everywhere
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, schemas.Error{
			Error: err.Error(),
		})
		return
	}
	s := `SELECT user_id, segments.segment_name, segments_history.action_date,  segments_history.action_type FROM segments_history
		JOIN segments USING(segment_id)
 		WHERE action_date BETWEEN $1 AND $2
 		ORDER BY user_id`
	rows, err := database.Query(s, (*request.Begin).UTC(), (*request.End).UTC())

	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.Error{
			Error: fmt.Sprintf("%s: %s", "DB error", err.Error()),
		})
		return
	}
	counter := 0
	for rows.Next() {
		userId := 0
		counter++
		n := ""
		var d time.Time
		act := ""
		rows.Scan(&userId, &n, &d, &act)
		err := w.Write([]string{strconv.Itoa(userId), n, act, d.Format(time.RFC822)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, schemas.Error{
				Error: fmt.Sprintf("%s: %s", "Server error", err.Error()),
			})
			return
		}
	}
	w.Flush()
	log.Print(counter)
	c.Header("Content-Type", "text/csv")
	log.Print(buf.String())
	_, err = c.Writer.Write(buf.Bytes())
	c.JSON(http.StatusInternalServerError, schemas.Error{
		Error: fmt.Sprintf("%s: %s", "Server error", err.Error()),
	})
	c.Writer.Flush()

}
