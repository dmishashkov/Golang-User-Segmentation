package schemas

import "time"

type Error struct {
	Error string `json:"error"`
}

type AddSegment struct {
	Name string     `json:"segment_name"`
	Time *time.Time `json:"delete_date"`
}
