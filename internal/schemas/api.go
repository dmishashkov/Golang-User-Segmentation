package schemas

import "time"

type Error struct {
	Error string `json:"error"`
}

type AddSlug struct {
	Name string     `json:"slug_name"`
	Time *time.Time `json:"delete_date"`
}
