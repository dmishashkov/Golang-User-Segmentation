package db

import (
	"database/sql"
	"fmt"
	"github.com/dmishashkov/avito_test_task_2023/config"
	"github.com/dmishashkov/avito_test_task_2023/internal/schemas"
	_ "github.com/lib/pq"
	"sync"
)

var GetSegments = ``
var DeleteSegment = ``
var CreateSegment = ``
var DeleteSegmentToUser = ``
var AddSegmentToUser = ``

func ConnectToDB(cfg schemas.DatabaseConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := sql.Open("postgres", connectionString)
	//if err != nil {
	//	log.Fatal("[ConnectToDB] error while connecting to DB", err)
	//}
	return db, err
}

var singleton sync.Once
var DB *sql.DB

func GetDB() (*sql.DB, error) {
	var Err error
	singleton.Do(func() {
		DB, Err = ConnectToDB(config.ProjectConfig.DB)
	})
	return DB, Err

}
