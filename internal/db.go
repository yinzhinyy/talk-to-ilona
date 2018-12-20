package internal

type DB interface {
	Save(db string, table string, document interface{}) interface{}
}

func LoadDB(dialect string) DB {
	var db DB
	switch dialect {
	case "mongo":
		db = &MongoDB{
			URL: "mongodb://35.247.89.169:27017",
		}
	}
	return db
}
