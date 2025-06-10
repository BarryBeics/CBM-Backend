package database

import (
	"context"
	"time"

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

	now := time.Now()

	user := &model.User{
		ID:             primitive.NewObjectID().Hex(),
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		Password:       string(hashedPassword),
		MobileNumber:   input.MobileNumber,
		Role:           input.Role,
		InvitedBy:      input.InvitedBy,
		VerifiedEmail:  false, // default false
		VerifiedMobile: false, // default false
		OpenToTrade:    false, // default false
		JoinedBallot:   false, // default false
		IsPaidMember:   false, // default false
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting user into the database:")
		return nil, err
	}

	return user, nil
}

// ReadUserByEmail retrieves a user by their email address.
func (db *DB) ReadUserByEmail(ctx context.Context, email string) (*model.User, error) {
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

// ReadUserByRole retrieves a group of users by their role.
func (db *DB) ReadUserByRole(ctx context.Context, role string) ([]*model.User, error) {
	collection := db.client.Database("go_trading_db").Collection("Customers")

	filter := bson.D{{"role", role}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Error querying users by role:")
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		log.Error().Err(err).Msg("Error decoding users by role:")
		return nil, err
	}

	return users, nil
}

// ReadAllUsers fetches all users from the Customers collection.
func (db *DB) ReadAllUsers(ctx context.Context) ([]*model.User, error) {
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

// UpdateUser modifies an existing user's details in the Customers collection.
func (db *DB) UpdateUser(ctx context.Context, input model.UpdateUserInput) (*model.User, error) {
	collection := db.client.Database("go_trading_db").Collection("Customers")

	filter := bson.M{"id": input.ID}
	updateFields := bson.M{
		"upatedAt": time.Now(), // your schema has a typo 'upatedAt'—you might want to fix it
	}

	if input.FirstName != nil {
		updateFields["firstName"] = *input.FirstName
	}
	if input.LastName != nil {
		updateFields["lastName"] = *input.LastName
	}
	if input.Email != nil {
		updateFields["email"] = *input.Email
	}
	if input.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error().Err(err).Msg("Error hashing updated password")
			return nil, err
		}
		updateFields["password"] = string(hashedPassword)
	}
	if input.MobileNumber != nil {
		updateFields["mobileNumber"] = *input.MobileNumber
	}
	if input.VerifiedEmail != nil {
		updateFields["verifiedEmail"] = *input.VerifiedEmail
	}
	if input.VerifiedMobile != nil {
		updateFields["verifiedMobile"] = *input.VerifiedMobile
	}
	if input.Role != nil {
		updateFields["role"] = *input.Role
	}
	if input.OpenToTrade != nil {
		updateFields["openToTrade"] = *input.OpenToTrade
	}
	if input.BinanceAPI != nil {
		updateFields["binanceAPI"] = *input.BinanceAPI
	}
	if input.PreferredContactMethod != nil {
		updateFields["preferredContactMethod"] = *input.PreferredContactMethod
	}
	if input.Notes != nil {
		updateFields["notes"] = *input.Notes
	}
	if input.InvitedBy != nil {
		updateFields["invitedBy"] = *input.InvitedBy
	}
	if input.JoinedBallot != nil {
		updateFields["joinedBallot"] = *input.JoinedBallot
	}
	if input.IsPaidMember != nil {
		updateFields["isPaidMember"] = *input.IsPaidMember
	}

	// This is required, because your schema expects it to be provided explicitly
	updateFields["isDeleted"] = input.IsDeleted

	update := bson.M{"$set": updateFields}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Error updating user:")
		return nil, err
	}

	var updated model.User
	err = collection.FindOne(ctx, filter).Decode(&updated)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving updated user:")
		return nil, err
	}

	return &updated, nil
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
