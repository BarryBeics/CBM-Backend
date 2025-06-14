package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
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

var seedDataDir = "/usr/local/share/seeds"

func init() {
	if custom := os.Getenv("SEED_DATA_DIR"); custom != "" {
		seedDataDir = custom
	}
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	_ = godotenv.Load()

	mode := flag.String("mode", "export", "Mode to run: export, import, or seed-users")
	flag.Parse()

	switch *mode {
	case "export":
		if err := runExport(); err != nil {
			log.Error().Err(err).Msg("Export failed")
		}
	case "import":
		if err := runImport(); err != nil {
			log.Error().Err(err).Msg("Import failed")
		}
	case "seed-users":
		if err := seedUsers(); err != nil {
			log.Error().Err(err).Msg("User seeding failed")
		}
	default:
		log.Warn().Str("mode", *mode).Msg("Unknown mode. Use 'export', 'import', or 'seed-users'")
	}
}

func runExport() error {
	client := graphql.NewClient("http://cbm-api:8080/query", http.DefaultClient)
	ctx := context.Background()

	log.Info().Msg("Fetching projects...")
	projectsResp, err := graph.ReadProjectsFilter(ctx, client, false)
	if err != nil {
		return fmt.Errorf("fetching projects: %w", err)
	}

	log.Info().Msg("Fetching SOPs...")
	sopsResp, err := graph.ReadProjectsFilter(ctx, client, true)
	if err != nil {
		return fmt.Errorf("fetching sops: %w", err)
	}

	log.Info().Msg("Fetching tasks...")
	tasksResp, err := graph.ReadAllTasks(ctx, client)
	if err != nil {
		return fmt.Errorf("fetching tasks: %w", err)
	}

	if err := writeJSON("projects.json", projectsResp.ReadProjectsFilter); err != nil {
		return fmt.Errorf("writing projects.json: %w", err)
	}
	log.Info().Msg("projects.json written")

	if err := writeJSON("sops.json", sopsResp.ReadProjectsFilter); err != nil {
		return fmt.Errorf("writing sops.json: %w", err)
	}
	log.Info().Msg("sops.json written")

	if err := writeJSON("tasks.json", tasksResp.ReadAllTasks); err != nil {
		return fmt.Errorf("writing tasks.json: %w", err)
	}
	log.Info().Msg("tasks.json written")

	log.Info().Msg("Export complete.")
	return nil
}

func runImport() error {
	var projects []ExportProject
	var sops []ExportProject
	var tasks []ExportTask

	log.Info().Msg("Reading JSON files...")

	if err := readJSON("projects.json", &projects); err != nil {
		return fmt.Errorf("reading projects.json: %w", err)
	}
	log.Info().Int("count", len(projects)).Msg("Loaded projects.json")

	if err := readJSON("sops.json", &sops); err != nil {
		return fmt.Errorf("reading sops.json: %w", err)
	}
	log.Info().Int("count", len(sops)).Msg("Loaded sops.json")

	if err := readJSON("tasks.json", &tasks); err != nil {
		return fmt.Errorf("reading tasks.json: %w", err)
	}
	log.Info().Int("count", len(tasks)).Msg("Loaded tasks.json")

	client := graphql.NewClient("http://cbm-api:8080/query", http.DefaultClient)
	ctx := context.Background()

	log.Info().Msg("Creating projects and SOPs...")
	for _, proj := range append(projects, sops...) {
		input := graph.CreateProjectInput{
			Title:       proj.Title,
			Sop:         proj.Sop,
			Description: proj.Description,
			Labels:      proj.Labels,
			AssignedTo:  proj.AssignedTo,
			DueDate:     proj.DueDate,
			Status:      proj.Status,
		}

		_, err := graph.CreateProject(ctx, client, input)
		if err != nil {
			log.Error().Str("project", proj.Title).Err(err).Msg("Failed to create project")
		} else {
			log.Info().Str("project", proj.Title).Msg("Project created")
		}
	}

	log.Info().Msg("Creating tasks...")
	for _, task := range tasks {
		input := graph.CreateTaskInput{
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			Labels:      []string{},
			AssignedTo:  task.AssignedTo,
			DueDate:     task.DueDate,
			ProjectId:   task.ProjectID,
		}

		_, err := graph.CreateTask(ctx, client, input)
		if err != nil {
			log.Error().Str("task", task.Title).Err(err).Msg("Failed to create task")
		} else {
			log.Info().Str("task", task.Title).Msg("Task created")
		}
	}

	log.Info().Msg("Import complete.")
	return nil
}

func writeJSON(filename string, data interface{}) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Error().Err(err).Str("file", filename).Msg("Failed to marshal JSON")
		return err
	}
	err = os.WriteFile(filename, bytes, 0644)
	if err != nil {
		log.Error().Err(err).Str("file", filename).Msg("Failed to write file")
	}
	return err
}

func readJSON(filename string, out interface{}) error {
	// Prepend default seed data path
	fullPath := fmt.Sprintf("%s/%s", seedDataDir, filename)

	bytes, err := os.ReadFile(fullPath)
	if err != nil {
		log.Error().Err(err).Str("file", fullPath).Msg("Failed to read file")
		return err
	}
	if err := json.Unmarshal(bytes, out); err != nil {
		log.Error().Err(err).Str("file", fullPath).Msg("Failed to unmarshal JSON")
		return err
	}
	return nil
}

func seedUsers() error {
	var users []model.CreateUserInput
	if err := readJSON("users.json", &users); err != nil {
		return fmt.Errorf("reading users.json: %w", err)
	}

	client := graphql.NewClient("http://cbm-api:8080/query", http.DefaultClient)
	ctx := context.Background()

	for _, u := range users {
		log.Info().Str("email", u.Email).Msg("Checking if user exists...")

		exists, err := userExists(ctx, client, u.Email)
		if err != nil {
			log.Error().Str("email", u.Email).Err(err).Msg("Error checking user existence")
			continue
		}
		if exists {
			log.Info().Str("email", u.Email).Msg("User already exists — skipping")
			continue
		}

		input := graph.CreateUserInput{
			Email:     u.Email,
			Password:  u.Password,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
		}

		_, err = graph.CreateUser(ctx, client, input)
		if err != nil {
			log.Error().Str("email", u.Email).Err(err).Msg("Failed to create user")
		} else {
			log.Info().Str("email", u.Email).Msg("User created")
		}
	}

	log.Info().Int("total", len(users)).Msg("User seeding complete")
	return nil
}

func userExists(ctx context.Context, client graphql.Client, email string) (bool, error) {
	resp, err := graph.ReadUserByEmail(ctx, client, email)
	if err != nil {
		// If the error is "no documents in result", treat as "user not found"
		if strings.Contains(err.Error(), "no documents in result") {
			return false, nil
		}
		return false, fmt.Errorf("GraphQL query failed: %w", err)
	}

	// Use the ID check as a fallback in case it returns an empty struct
	return resp.ReadUserByEmail.Id != "", nil
}
