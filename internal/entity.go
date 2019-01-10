package internal

import "time"

type MessageBody struct {
	UserName   string
	Year       int
	Month      time.Month
	Day        int
	CreateTime int64
	Activity   string
}
