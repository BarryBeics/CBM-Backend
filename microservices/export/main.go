package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"cryptobotmanager.com/cbm-backend/shared/graph"
	"github.com/Khan/genqlient/graphql"
)

type ExportProject struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Sop         bool     `json:"sop"`
	Description string   `json:"description"`
	Labels      []string `json:"labels"`
	AssignedTo  string   `json:"assignedTo"`
	DueDate     string   `json:"dueDate"`
	Status      string   `json:"status"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

type ExportTask struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	AssignedTo  string `json:"assignedTo"`
	DueDate     string `json:"dueDate"`
	Category    string `json:"category"`
	ProjectID   string `json:"projectId"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func main() {
	client := graphql.NewClient("http://cbm-api:8080/query", http.DefaultClient) // Change to your actual endpoint
	ctx := context.Background()

	projects, err := graph.GetAllProjects(ctx, client, false)
	if err != nil {
		fmt.Printf("Error fetching projects: %v\n", err)
		os.Exit(1)
	}

	sops, err := graph.GetAllProjects(ctx, client, true)
	if err != nil {
		fmt.Printf("Error fetching projects: %v\n", err)
		os.Exit(1)
	}

	tasks, err := graph.GetAllTasks(ctx, client)
	if err != nil {
		fmt.Printf("Error fetching tasks: %v\n", err)
		os.Exit(1)
	}

	if err := writeJSON("projects.json", projects); err != nil {
		fmt.Printf("Error writing projects.json: %v\n", err)
	}

	if err := writeJSON("sops.json", sops); err != nil {
		fmt.Printf("Error writing sops.json: %v\n", err)
	}

	if err := writeJSON("tasks.json", tasks); err != nil {
		fmt.Printf("Error writing tasks.json: %v\n", err)
	}

	fmt.Println("✅ Export complete.")
}

func writeJSON(filename string, data interface{}) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, bytes, 0644)
}
