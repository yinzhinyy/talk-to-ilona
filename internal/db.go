package internal

import (
	"log"
	"os"
)

type DB interface {
	Save(db string, table string, document interface{}) interface{}
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
