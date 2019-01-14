package internal

import "time"

// MessageBody database entity
type MessageBody struct {
	UserName   string
	Year       int
	Month      time.Month
	Day        int
	CreateTime int64
	Activity   string
}
