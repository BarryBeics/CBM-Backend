package database

import (
	"context"

	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateTask inserts a new task into the Tasks collection.
func (db *DB) CreateTask(ctx context.Context, input model.CreateTaskInput) (*model.Task, error) {
	collection := db.client.Database("go_trading_db").Collection("Tasks")

	now := time.Now().Format(time.RFC3339)

	task := bson.M{
		"id":          primitive.NewObjectID().Hex(),
		"title":       input.Title,
		"description": input.Description,
		"status":      *input.Status,
		"labels":      input.Labels,
		"assignedTo":  input.AssignedTo,
		"dueDate":     input.DueDate,
		"deferDate":   input.DeferDate,
		"department":  input.Department,
		"projectId":   input.ProjectID, // ensure this stays a string
		"duration":    input.Duration,
		"createdAt":   now,
		"updatedAt":   now,
	}

	var updated model.Task
	_, err := collection.InsertOne(ctx, task)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting task into the database:")
		return nil, err
	}

	return &updated, nil
}

// UpdateTask updates an existing task in the Tasks collection.
func (db *DB) UpdateTask(ctx context.Context, input model.UpdateTaskInput) (*model.Task, error) {
	collection := db.client.Database("go_trading_db").Collection("Tasks")

	filter := bson.M{"id": input.ID}
	updateFields := bson.M{
		"updatedAt": time.Now().Format(time.RFC3339),
	}

	// Only update fields that are non-nil (optional fields)
	if input.Title != nil {
		updateFields["title"] = input.Title
	}
	if input.Description != nil {
		updateFields["description"] = input.Description
	}
	if input.Status != nil {
		updateFields["status"] = input.Status
	}
	if input.Labels != nil {
		updateFields["labels"] = input.Labels
	}
	if input.AssignedTo != nil {
		updateFields["assignedTo"] = input.AssignedTo
	}
	if input.DueDate != nil {
		updateFields["dueDate"] = input.DueDate
	}
	if input.DeferDate != nil {
		updateFields["deferDate"] = input.DeferDate
	}
	if input.Department != nil {
		updateFields["department"] = input.Department
	}
	if input.ProjectID != nil {
		updateFields["projectId"] = input.ProjectID
	}
	if input.Duration != nil {
		updateFields["duration"] = input.Duration
	}

	update := bson.M{"$set": updateFields}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Error updating task:")
		return nil, err
	}

	var updated model.Task
	err = collection.FindOne(ctx, filter).Decode(&updated)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving updated task:")
		return nil, err
	}

	return &updated, nil
}

// DeleteTaskByID removes a task by its ID.
func (db *DB) DeleteTaskByID(ctx context.Context, id string) (bool, error) {
	collection := db.client.Database("go_trading_db").Collection("Tasks")

	filter := bson.M{"id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting task:")
		return false, err
	}

	return result.DeletedCount > 0, nil
}

// GetAllTasks fetches all tasks from the Tasks collection.
func (db *DB) GetAllTasks(ctx context.Context) ([]*model.Task, error) {
	collection := db.client.Database("go_trading_db").Collection("Tasks")

	cursor, err := collection.Find(ctx, bson.D{}, options.Find())
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving tasks:")
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*model.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		log.Error().Err(err).Msg("Error decoding tasks:")
		return nil, err
	}

	return tasks, nil
}

// GetTaskByID fetches a single task by its ID.
func (db *DB) GetTaskByID(ctx context.Context, id string) (*model.Task, error) {
	collection := db.client.Database("go_trading_db").Collection("Tasks")

	filter := bson.M{"id": id}

	var task model.Task
	err := collection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving task by ID:")
		return nil, err
	}

	return &task, nil
}

// CreateProject inserts a new project into the Projects collection.
func (db *DB) CreateProject(ctx context.Context, input model.CreateProjectInput) (*model.Project, error) {
	collection := db.client.Database("go_trading_db").Collection("Projects")

	now := time.Now().Format(time.RFC3339)

	project := &model.Project{
		ID:          primitive.NewObjectID().Hex(),
		Title:       input.Title,
		Sop:         *input.Sop,
		Description: input.Description,
		Labels:      input.Labels,
		AssignedTo:  input.AssignedTo,
		DueDate:     input.DueDate,
		Status:      *input.Status,

		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := collection.InsertOne(ctx, project)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting project into the database:")
		return nil, err
	}

	return project, nil
}

// UpdateProject updates an existing project in the Projects collection.
func (db *DB) UpdateProject(ctx context.Context, input model.UpdateProjectInput) (*model.Project, error) {
	collection := db.client.Database("go_trading_db").Collection("Projects")

	filter := bson.M{"id": input.ID}
	updateFields := bson.M{
		"updatedAt": time.Now().Format(time.RFC3339),
	}

	if input.Title != nil {
		updateFields["title"] = input.Title
	}
	if input.Description != nil {
		updateFields["description"] = input.Description
	}
	if input.AssignedTo != nil {
		updateFields["assignedTo"] = input.AssignedTo
	}
	if input.DueDate != nil {
		updateFields["dueDate"] = input.DueDate
	}
	if input.Status != nil {
		updateFields["status"] = input.Status
	}
	if input.Labels != nil {
		updateFields["labels"] = input.Labels
	}
	if input.Sop != nil {
		updateFields["sop"] = input.Sop
	}

	update := bson.M{"$set": updateFields}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Error updating project:")
		return nil, err
	}

	var updated model.Project
	err = collection.FindOne(ctx, filter).Decode(&updated)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving updated project:")
		return nil, err
	}

	return &updated, nil
}

// DeleteProjectByID deletes a project by ID.
func (db *DB) DeleteProjectByID(ctx context.Context, id string) (bool, error) {
	collection := db.client.Database("go_trading_db").Collection("Projects")

	filter := bson.M{"id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting project:")
		return false, err
	}

	return result.DeletedCount > 0, nil
}

// GetProjectsBySOP retrieves projects filtered by the SOP boolean field.
func (db *DB) GetProjectsBySOP(ctx context.Context, sop bool) ([]*model.Project, error) {
	collection := db.client.Database("go_trading_db").Collection("Projects")

	// Build the MongoDB filter
	filter := bson.M{"sop": sop}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving filtered projects:")
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []*model.Project
	if err := cursor.All(ctx, &projects); err != nil {
		log.Error().Err(err).Msg("Error decoding filtered projects:")
		return nil, err
	}

	return projects, nil
}

// GetProjectByID retrieves a project by ID.
func (db *DB) GetProjectByID(ctx context.Context, id string) (*model.Project, error) {
	collection := db.client.Database("go_trading_db").Collection("Projects")

	filter := bson.M{"id": id}

	var project model.Project
	err := collection.FindOne(ctx, filter).Decode(&project)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving project by ID:")
		return nil, err
	}

	return &project, nil
}

// GetTasksByProjectID retrieves tasks associated with a specific project.
func (db *DB) GetTasksByProjectID(ctx context.Context, projectID string) ([]*model.Task, error) {
	collection := db.client.Database("go_trading_db").Collection("Tasks")

	filter := bson.M{"projectId": projectID}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving tasks by project ID:")
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*model.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		log.Error().Err(err).Msg("Error decoding tasks by project ID:")
		return nil, err
	}

	return tasks, nil
}
