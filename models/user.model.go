package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserImage struct {
	PublicID string `json:"public_id" bson:"public_id"`
	URL      string `json:"url" bson:"url"`
}

type USER struct {
	Id                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email              string             `json:"email" bson:"email" validate:"required,email"`
	FirstName          string             `json:"first_name" bson:"first_name,min=2,max=100" `
	UserName           string             `json:"user_name" bson:"user_name" validate:"required,min=2,max=100"`
	LastName           string             `json:"last_name" bson:"last_name" validate:"min=2,max=100"`
	Phone_Number       string             `json:"phone_number" bson:"phone_number" validate:"min=2,max=100"`
	Password           string             `json:"password" bson:"password" validate:"required"`
	UserRole           string             `json:"user_role" bson:"user_role" validate:"required, eq=ADMIN|eq=USER"`
	Created_At         time.Time          `json:"created_at" bson:"created_at"`
	Updated_At         time.Time          `json:"updated_at" bson:"updated_at"`
	ResetPasswordToken string             `json:"reset_password_token" bson:"reset_password_token"`
	UserImage          UserImage          `json:"user_image" bson:"user_image"`
}
