package model

import (
	"io"
	"time"
)

type Attachment struct {
	Name        string
	ContentType string
	Data        io.Reader
}

type Message struct {
	Id             string
	Text           string
	User           string
	DoNotPrintUser bool
	Time           time.Time
	Attachments    []Attachment
}
