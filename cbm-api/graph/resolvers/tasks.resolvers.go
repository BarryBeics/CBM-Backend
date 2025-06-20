package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.72

import (
	"context"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"github.com/rs/zerolog/log"
)

// CreateTask is the resolver for the createTask field.
func (r *mutationResolver) CreateTask(ctx context.Context, input model.CreateTaskInput) (*model.Task, error) {
	task, err := db.CreateTask(ctx, input)
	if err != nil {
		log.Error().Err(err).Msg("Error creating task:")
		return nil, err
	}

	return task, nil
}

// UpdateTask is the resolver for the updateTask field.
func (r *mutationResolver) UpdateTask(ctx context.Context, input model.UpdateTaskInput) (*model.Task, error) {
	task, err := db.UpdateTask(ctx, input)
	if err != nil {
		log.Error().Err(err).Msg("Error updating task:")
		return nil, err
	}

	return task, nil
}

// DeleteTask is the resolver for the deleteTask field.
func (r *mutationResolver) DeleteTask(ctx context.Context, id string) (*bool, error) {
	success, err := db.DeleteTaskByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting task:")
		return nil, err
	}

	return &success, nil
}

// CreateProject is the resolver for the createProject field.
func (r *mutationResolver) CreateProject(ctx context.Context, input model.CreateProjectInput) (*model.Project, error) {
	project, err := db.CreateProject(ctx, input)
	if err != nil {
		log.Error().Err(err).Msg("Error creating project:")
		return nil, err
	}
	return project, nil
}

// UpdateProject is the resolver for the updateProject field.
func (r *mutationResolver) UpdateProject(ctx context.Context, input model.UpdateProjectInput) (*model.Project, error) {
	project, err := db.UpdateProject(ctx, input)
	if err != nil {
		log.Error().Err(err).Msg("Error updating project:")
		return nil, err
	}
	return project, nil
}

// DeleteProject is the resolver for the deleteProject field.
func (r *mutationResolver) DeleteProject(ctx context.Context, id string) (*bool, error) {
	success, err := db.DeleteProjectByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting project:")
		return nil, err
	}
	return &success, nil
}

// ReadTaskByID is the resolver for the readTaskById field.
func (r *queryResolver) ReadTaskByID(ctx context.Context, id string) (*model.Task, error) {
	task, err := db.ReadTaskByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching task by ID:")
		return nil, err
	}

	return task, nil
}

// ReadAllTasks is the resolver for the readAllTasks field.
func (r *queryResolver) ReadAllTasks(ctx context.Context) ([]*model.Task, error) {
	tasks, err := db.ReadAllTasks(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching tasks:")
		return nil, err
	}

	return tasks, nil
}

// ReadSingleProjectByID is the resolver for the readSingleProjectById field.
func (r *queryResolver) ReadSingleProjectByID(ctx context.Context, id string) (*model.Project, error) {
	project, err := db.ReadSingleProjectByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching project by ID:")
		return nil, err
	}

	// Fetch tasks manually
	tasks, err := db.GetTasksByProjectID(ctx, project.ID)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching tasks for project:")
		return nil, err
	}
	project.Tasks = tasks

	return project, nil
}

// ReadProjectsFilter is the resolver for the readProjectsFilter field.
func (r *queryResolver) ReadProjectsFilter(ctx context.Context, filter *model.ProjectFilterInput) ([]*model.Project, error) {
	var projects []*model.Project
	var err error

	if filter != nil && filter.Sop != nil {
		projects, err = db.ReadProjectsBySOP(ctx, *filter.Sop)
	}

	if err != nil {
		log.Error().Err(err).Msg("Error fetching projects:")
		return nil, err
	}

	for _, project := range projects {
		tasks, err := db.GetTasksByProjectID(ctx, project.ID)
		if err != nil {
			log.Error().Err(err).Str("projectID", project.ID).Msg("Error fetching tasks")
			continue
		}
		project.Tasks = tasks
	}

	return projects, nil
}
