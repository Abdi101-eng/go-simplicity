package taskstore

import (
	"abdi/task-manager/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileStore struct {
	path string
	data fileData
}

type fileData struct {
	Tasks  []models.Task `json:"tasks"`
	NextID int           `json:"next_id"`
}

const (
	defaultDirName  = ".taskd"
	defaultFileName = "tasks.json"
)

func NewFileStore(path string) (*FileStore, error) {
	if path == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		path = filepath.Join(homedir, defaultDirName, defaultFileName)
	}

	file, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &FileStore{
				path: path,
			}, nil
		}
		return nil, err
	}

	var filedata fileData
	err = json.Unmarshal(file, &filedata)
	if err != nil {
		return nil, fmt.Errorf("filestore: %w", err)
	}

	return &FileStore{path: path, data: filedata}, nil

}

func (f *FileStore) Save() error {
	dir := filepath.Dir(f.path)
	// create the dir if doesnt already exist
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	jsonEncoding, err := json.Marshal(f.data)
	if err != nil {
		return err
	}

	file, err := os.CreateTemp(dir, "tasks-*.tmp")
	if err != nil {
		return err
	}
	defer file.Close()

	err = os.WriteFile(file.Name(), jsonEncoding, 0755)
	if err != nil {
		return err
	}

	err = os.Rename(file.Name(), f.path)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileStore) Add(title string, priority models.Priority) (models.Task, error) {
	// When i want to add a task, i want to append to the fileData and then click save?
	f.data.NextID++
	task := models.Task{
		Id:        f.data.NextID,
		Title:     title,
		Priority:  priority,
		Done:      false,
		CreatedAt: time.Now(),
	}
	f.data.Tasks = append(f.data.Tasks, task)
	err := f.Save()
	if err != nil {
		return models.Task{}, fmt.Errorf("filestore: %w", err)
	}
	return task, nil
}

func (f *FileStore) List() []models.Task {
	cp := make([]models.Task, len(f.data.Tasks))
	copy(cp, f.data.Tasks)
	return cp
}

func (f *FileStore) Complete(id int) error {
	for index, element := range f.data.Tasks {
		if element.Id == id {
			f.data.Tasks[index].Done = true
			err := f.Save()
			if err != nil {
				return fmt.Errorf("complete: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("complete: %w", models.ErrTaskNotFound)
}

func (f *FileStore) Delete(id int) error {
	for index, element := range f.data.Tasks {
		if element.Id == id {
			f.data.Tasks = append(f.data.Tasks[:index], f.data.Tasks[index+1:]...)
			err := f.Save()
			if err != nil {
				return fmt.Errorf("delete: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("delete: %w", models.ErrTaskNotFound)
}
