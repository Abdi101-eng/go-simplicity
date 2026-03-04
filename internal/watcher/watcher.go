package watcher

import (
	"abdi/task-manager/internal/models"
	"context"
	"fmt"
	"time"
)

func StartWatcher(ctx context.Context, store models.TaskStore, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			checkOverdueTasks(store)
		case <-ctx.Done():
			return
		}
	}
}

func checkOverdueTasks(store models.TaskStore) {
	tasks := store.List()
	pending := make([]models.Task, 0)
	for _, task := range tasks {
		if !task.Done {
			pending = append(pending, task)
		}
	}
	if len(pending) == 0 {
		return
	}
	fmt.Println("Pending tasks:")
	for _, task := range pending {
		fmt.Printf("  - %s [%s]\n", task.Title, task.Priority)
	}
}
