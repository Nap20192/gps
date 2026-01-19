package models

import "github.com/google/uuid"

type User struct {
	UserID       uuid.UUID `json:"id" bson:"id"`
	Username     string    `json:"username" bson:"username"`
	PasswordHash string    `json:"password" bson:"password"`
}
