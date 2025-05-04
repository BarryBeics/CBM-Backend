package database

import (
	"context"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser inserts a new user into the Customers collection.
func (db *DB) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	collection := db.client.Database("go_trading_db").Collection("Customers")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Error hashing password")
		return nil, err
	}

	user := &model.User{
		ID:        primitive.NewObjectID().Hex(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Contact:   input.Contact,
		Address1:  input.Address1,
		Address2:  input.Address2,
		Role:      input.Role,
		Password:  string(hashedPassword),
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting user into the database:")
		return nil, err
	}

	return user, nil
}

// DeleteUserByEmail removes a user from the database using their email.
func (db *DB) DeleteUserByEmail(ctx context.Context, email string) (bool, error) {
	collection := db.client.Database("go_trading_db").Collection("Customers")

	filter := bson.D{{"email", email}}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting user from the database:")
		return false, err
	}

	return result.DeletedCount > 0, nil
}

// GetUserByEmail retrieves a user by their email address.
func (db *DB) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	collection := db.client.Database("go_trading_db").Collection("Customers")

	filter := bson.D{{"email", email}}

	var user model.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving user by email:")
		return nil, err
	}

	return &user, nil
}

// GetAllUsers fetches all users from the Customers collection.
func (db *DB) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	collection := db.client.Database("go_trading_db").Collection("Customers")

	cursor, err := collection.Find(ctx, bson.D{}, options.Find())
	if err != nil {
		log.Error().Err(err).Msg("Error querying all users:")
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		log.Error().Err(err).Msg("Error decoding users:")
		return nil, err
	}

	return users, nil
}
