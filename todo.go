package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	Time      time.Time `json:"time"`
}

type TodoList struct {
	Tasks    []Task `json:"tasks"`
	FilePath string `json:"filePath"`
	NextId   int    `json:"nextId"`
}

func NewTodoList(FilePath string) *TodoList {
	list := &TodoList{
		Tasks:    []Task{},
		FilePath: FilePath,
		NextId:   1,
	}
	list.Load()
	return list
}

func (tl *TodoList) Load() error {
	_, err := os.Stat(tl.FilePath)
	if os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(tl.FilePath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}
	err = json.Unmarshal(data, tl)
	return err
}

func (tl *TodoList) AddTask(title string) Task {
	task := Task{
		ID:        tl.NextId,
		Title:     title,
		Completed: false,
		Time:      time.Now(),
	}
	tl.Tasks = append(tl.Tasks, task)
	tl.NextId += 1
	tl.Save()
	return task
}

func (tl *TodoList) ListTasks(showCompleted bool) []Task {
	if showCompleted {
		return tl.Tasks
	}
	var activeTasks []Task
	for _, task := range tl.Tasks {
		if !task.Completed {
			activeTasks = append(activeTasks, task)
		}
	}
	return activeTasks
}

func (tl *TodoList) CompleteTask(taskId int) bool {
	for i, task := range tl.Tasks {
		if task.ID == taskId {
			tl.Tasks[i].Completed = true
			tl.Save()
			return true
		}
	}
	return false
}

func (tl *TodoList) DeleteTask(id int) bool {
	for i, task := range tl.Tasks {
		if task.ID == id {
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
			tl.Save()
			return true
		}
	}
	return false
}

func (tl *TodoList) Save() error {
	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tl.FilePath, data, 0644)
}

func main() {
	filePath := "./todo.json"
	todoList := NewTodoList(filePath)
	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing task title")
			return
		}
		task := todoList.AddTask(os.Args[2])
		fmt.Printf("Added task %s\n", task.Title)
	case "list":
		showCompleted := false
		if len(os.Args) > 2 && os.Args[2] == "--all" {
			showCompleted = true
		}

		tasks := todoList.ListTasks(showCompleted)
		if len(tasks) == 0 {
			fmt.Println("No tasks available")
			return
		}
		fmt.Println(" ID | Status | Title")
		fmt.Println("-------------------------------------")
		for _, task := range tasks {
			status := "[ ]"
			if task.Completed {
				status = "[✓]"
			}
			fmt.Printf(" %2d |   %s  | %s\n", task.ID, status, task.Title)
		}
	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing task ID")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: Invalid task ID")
			return
		}
		if todoList.CompleteTask(id) {
			fmt.Printf("Marked task %d as completed\n", id)
		} else {
			fmt.Printf("Task %d not found\n", id)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing task ID")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: Invalid task ID")
			return
		}
		if todoList.DeleteTask(id) {
			fmt.Printf("Deleted task %d\n", id)
		} else {
			fmt.Printf("Task %d not found\n", id)
		}

	case "help":
		showAllCommands()

	default:
		fmt.Printf("%s is not a valid command. See 'todo help'", command)
	}
}

func showAllCommands() {
	fmt.Println("Commands")
	fmt.Println("  todo add <task title>     Add a new task")
	fmt.Println("  todo list                 List active tasks")
	fmt.Println("  todo list --all           List all tasks including completed")
	fmt.Println("  todo complete <id>        Mark a task as completed")
	fmt.Println("  todo delete <id>          Delete a task")
	fmt.Println("  todo help                 Show available commands")
}
