package internal

import (
	"log"
	"os"

	"github.com/mongodb/mongo-go-driver/bson"
)

const (
	// DBName is database namespace
	DBName = "talk"
	// DBTableActivity is table name of activity log
	DBTableActivity = "activity_log"
	// DBColumnUserName is column userName of table activity log
	DBColumnUserName = "userName"
	// DBColumnYear is column year (type: int) of table activity log
	DBColumnYear = "year"
	// DBColumnMonth is column month (type: int) of table activity log
	DBColumnMonth = "month"
	// DBColumnDay is column day (type: int) of table activity log
	DBColumnDay = "day"
	// DBColumnActivities is column activities (list of activity) of table activity log
	DBColumnActivities = "activities"
	// DBColumnCreateTime is column createTime (type: int64) of table activity log
	DBColumnCreateTime = "createTime"
	// DBColumnActivity is column activity of table activity log
	DBColumnActivity = "activity"
)

// DB is super class of database api
type DB interface {
	Save(db string, table string, document interface{}) interface{}
	Find(db string, table string, filter interface{}) bson.M
	Update(db string, table string, filter interface{}, document interface{}) interface{}
}

// LoadDB init database connection
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
