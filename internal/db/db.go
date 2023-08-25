package db

import (
	"database/sql"
	"fmt"
	"github.com/dmishashkov/avito_test_task_2023/config"
	"github.com/dmishashkov/avito_test_task_2023/internal/schemas"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

var GetSegments = ``
var DeleteSegment = ``
var CreateSegment = ``
var DeleteSegmentToUser = ``
var AddSegmentToUser = ``

func ConnectToDB(cfg schemas.DatabaseConfig) *sql.DB {
	connectionString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("[ConnectToDB] error while connecting to DB", err)
	}
	log.Println("Successfully connected to DB")
	return db
}

var singleton sync.Once
var myDB *sql.DB

func GetDB() *sql.DB {
	singleton.Do(func() {
		myDB = ConnectToDB(config.ProjectConfig.DB)
	})
	return myDB
}
