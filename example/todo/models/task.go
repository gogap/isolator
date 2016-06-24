package models

import (
	"time"
)

type Task struct {
	ID          string
	OwnerID     string
	Title       string
	Description string
	CreateTime  time.Time
	UpdateTime  time.Time
	Version     int
}
