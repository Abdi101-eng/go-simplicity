package models

import (
	"errors"
	"strings"
	"time"
)

type Task struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	Priority  Priority  `json:"priority"`
	CreatedAt time.Time `json:"createdat"`
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
	switch strings.ToLower(s) {
	case "low":
		return Low
	case "medium":
		return Medium
	case "high":
		return High
	default:
		return Low // Default to Low if the string doesn't match any known priority
	}
}

var ErrTaskNotFound = errors.New("task not found")

type TaskStore interface {
	Add(title string, priority Priority) (Task, error)
	List() []Task
	Complete(id int) error
	Delete(id int) error
}
