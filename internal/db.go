package internal

import (
	"log"
	"os"

	"github.com/mongodb/mongo-go-driver/bson"
)

const (
	TimeLayout         = "2006-01-02"
	DBName             = "talk"
	DBTableActivity    = "activity_log"
	DBColumnUserName   = "userName"
	DBColumnYear       = "year"
	DBColumnMonth      = "month"
	DBColumnDay        = "day"
	DBColumnCreateTime = "createTime"
	DBColumnActivities = "activities"
	DBColumnActivity   = "activity"
)

type DB interface {
	Save(db string, table string, document interface{}) interface{}
	Find(db string, table string, filter interface{}) bson.M
	Update(db string, table string, filter interface{}, document interface{}) interface{}
}

func LoadDB() DB {
	var dialect = os.Getenv("TALK_DB_TYPE")
	log.Printf("using env: TALK_DB_TYPE=%s", dialect)
	var db DB
	switch dialect {
	case "mongo":
		var url = os.Getenv("TALK_DB_MONGO_URL")
		log.Printf("using env: TALK_DB_MONGO_URL=%s", url)
		db = &MongoDB{
			URL: url,
		}
	}
	return db
}
