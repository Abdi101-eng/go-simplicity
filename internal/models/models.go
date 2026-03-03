package models

import (
	"errors"
	"time"
)

type Task struct {
	Id        int
	Title     string
	Done      bool
	Priority  Priority
	CreatedAt time.Time
}

type Priority int

const (
	Low Priority = iota + 1
	Medium
	High
)

func (P Priority) String() string {
	switch P {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	default:
		return "Unknown"
	}
}

func PriorityFromString(s string) Priority {
	switch s {
	case "Low":
		return Low
	case "Medium":
		return Medium
	case "High":
		return High
	default:
		return Low // Default to Low if the string doesn't match any known priority
	}
}

var ErrTaskNotFound = errors.New("task not found")

type TaskStore interface {
	Add(title string, priority Priority) Task
	List() []Task
	Complete(id int) error
	Delete(id int) error
}
