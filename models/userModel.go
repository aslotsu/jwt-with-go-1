package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `json:"ID,omitempty"`
	FirstName    string             `json:"firstName" validate:"required,min=2,max=100"`
	LastName     string             `json:"lastName" validate:"required,min=2,max=100"`
	Password     string             `json:"password" validate:"required,min=6"`
	Email        string             `json:"email" validate:"email,required"`
	PhoneNumber  string             `json:"phoneNumber" validate:"required"`
	Token        string             `json:"token"`
	RefreshToken string             `json:"refreshToken"`
	UserType     string             `json:"userType" validate:"required,eq=ADMIN|eq=USER"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
	UserId       string             `json:"userId"`
}
