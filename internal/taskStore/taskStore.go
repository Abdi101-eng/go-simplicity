package taskstore

import (
	"abdi/task-manager/internal/models"
	"fmt"
	"time"
)

var _ models.TaskStore = (*Store)(nil)

type Store struct {
	list   []models.Task
	nextId int
}

func NewStore() *Store {
	return &Store{
		list: make([]models.Task, 0),
	}
}

func (s *Store) Add(title string, priority models.Priority) models.Task {
	s.nextId++
	task := models.Task{
		Id:        s.nextId,
		Title:     title,
		Priority:  priority,
		Done:      false,
		CreatedAt: time.Now(),
	}
	s.list = append(s.list, task)
	return task
}

func (s *Store) List() []models.Task {
	out := make([]models.Task, len(s.list))
	copy(out, s.list)
	return out
}

func (s *Store) Complete(id int) error {
	for index, element := range s.list {
		if element.Id == id {
			s.list[index].Done = true
			return nil
		}
	}
	return fmt.Errorf("complete: %w", models.ErrTaskNotFound)
}

func (s *Store) Delete(id int) error {
	for index, element := range s.list {
		if element.Id == id {
			s.list = append(s.list[:index], s.list[index+1:]...)
			return nil
		}
	}
	return fmt.Errorf("delete: %w", models.ErrTaskNotFound)
}
