package main

import (
	"abdi/task-manager/internal/models"
	taskstore "abdi/task-manager/internal/taskStore"
	"abdi/task-manager/internal/watcher"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "no command provided. use 'help' for usage\n")
		os.Exit(1)
	}
	store, err := taskstore.NewFileStore("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	switch args[0] {
	case "add":
		err := runAddCommand(store, args)
		if err != nil {
			os.Exit(1)
		}
	case "list":
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			watcher.StartWatcher(ctx, store, 2*time.Second)
		}()

		runListCommand(store)
		<-ctx.Done()
		fmt.Println("\nshutting down...")
		wg.Wait()
	case "done":
		err := runDoneCommand(store, args)
		if err != nil {
			os.Exit(1)
		}
	case "delete":
		err := runDeleteCommand(store, args)
		if err != nil {
			os.Exit(1)
		}
	case "help":
		fmt.Println("Available commands: add, list, complete, delete")
		fmt.Println("Usage:")
		fmt.Println("  add <title> <priority> - Add a new task with the given title and priority (Low, Medium, High)")
		fmt.Println("  list - List all tasks")
		fmt.Println("  complete <id> - Mark the task with the given id as completed")
		fmt.Println("  delete <id> - Delete the task with the given id")
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s", args[0])
		os.Exit(1)
	}

}

func runAddCommand(store models.TaskStore, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Incorrect argument count")
	}
	title := args[1]
	priority := parsePriority(args)
	fmt.Printf("Adding task: %s with priority %s", title, priority)
	task, err := store.Add(title, models.PriorityFromString(priority))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}
	fmt.Printf("added task #%d\n", task.Id)
	return nil

}

func runListCommand(store models.TaskStore) {
	fmt.Printf("%-4s %-20s %-10s %s\n", "ID", "TITLE", "PRIORITY", "STATUS")
	for _, task := range store.List() {
		status := "pending"
		if task.Done {
			status = "done"
		}
		fmt.Printf("%-4d %-20s %-10s %s\n", task.Id, task.Title, task.Priority, status)
	}
}

func runDoneCommand(store models.TaskStore, args []string) error {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: done <id>")
		return fmt.Errorf("Usage: done <id>")
	}
	id := args[1]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid id %q\n", id)
		return err
	}
	err = store.Complete(idInt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}
	fmt.Printf("task #%d marked as done\n", idInt)
	return nil
}

func runDeleteCommand(store models.TaskStore, args []string) error {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: delete <id>")
		return fmt.Errorf("Usage: delete <id>")
	}
	id := args[1]
	idint, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid id %q\n", id)
		return err
	}
	err = store.Delete(idint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}
	fmt.Printf("task #%d deleted\n", idint)

	return nil
}

func parsePriority(args []string) string {
	for i, arg := range args {
		if arg == "--priority" && i+1 < len(args) {
			priority := args[i+1]
			return priority
		}
	}
	return "Medium" // Default priority
}
