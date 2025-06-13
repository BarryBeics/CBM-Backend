package database

import (
	"context"
	"os"
	"time"

	"cryptobotmanager.com/cbm-backend/cbm-api/graph/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// jwtSecret is the secret key used to sign JWT tokens.
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Login authenticates a user and returns a JWT token if successful.
func (db *DB) Login(ctx context.Context, input *model.LoginInput) (*model.LoginResponse, error) {
	collection := db.client.Database("go_trading_db").Collection("Customers")

	var user model.User
	err := collection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		log.Error().Err(err).Msg("User not found")
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		log.Error().Err(err).Msg("Invalid password")
		return nil, err
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"role":   user.Role,
		"exp":    time.Now().Add(time.Hour * 72).Unix(), // 3 days expiry
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Error().Err(err).Msg("Failed to sign JWT")
		return nil, err
	}

	return &model.LoginResponse{
		Token: tokenString,
		User:  &user,
	}, nil
}
