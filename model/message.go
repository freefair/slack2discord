package model

import "time"

type Message struct {
	Id   string
	Text string
	User string
	Time time.Time
}
