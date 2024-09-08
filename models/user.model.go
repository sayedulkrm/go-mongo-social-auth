package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// USER struct definition
type USER struct {
	Id                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email              string             `json:"email,omitempty" bson:"email,omitempty" validate:"required,email"`
	FirstName          *string            `json:"first_name,omitempty" bson:"first_name,omitempty" validate:"required"`
	LastName           *string            `json:"last_name,omitempty" bson:"last_name,omitempty" validate:"required,min=2,max=100"`
	Phone_Number       string             `json:"phone_number,omitempty" bson:"phone_number,omitempty" validate:"required,min=2,max=100"`
	Password           string             `json:"password,omitempty" bson:"password,omitempty" validate:"required"`
	UserRole           string             `json:"user_role,omitempty" validate:"required, eq=ADMIN|eq=USER"`
	Created_At         time.Time          `json:"created_at"`
	Updated_At         time.Time          `json:"updated_at"`
	ResetPasswordToken string             `json:"reset_password_token,omitempty"`
}
