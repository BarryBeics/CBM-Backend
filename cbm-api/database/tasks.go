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

	task := &model.Task{
		ID:          primitive.NewObjectID().Hex(),
		Title:       input.Title,
		Description: input.Description,
		Category:    input.Category,
		AssignedTo:  input.AssignedTo,
		DueDate:     input.DueDate,
		SopLink:     input.SopLink,
		Priority:    input.Priority,
		Status:      *input.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err := collection.InsertOne(ctx, task)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting task into the database:")
		return nil, err
	}

	return task, nil
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
	if input.Category != nil {
		updateFields["category"] = input.Category
	}
	if input.AssignedTo != nil {
		updateFields["assignedTo"] = input.AssignedTo
	}
	if input.DueDate != nil {
		updateFields["dueDate"] = input.DueDate
	}
	if input.SopLink != nil {
		updateFields["sopLink"] = input.SopLink
	}
	if input.Priority != nil {
		updateFields["priority"] = input.Priority
	}
	if input.Status != nil {
		updateFields["status"] = input.Status
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
