package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)
type Address struct {
	District *string `json:"district" validation:"required"`
	Parish   *string `json:"parish" validation:"required"`
	Village  *string `json:"village" bson:"village"`
}

type User struct {
	ID      primitive.ObjectID `bson:"_id"`
	FirstName    *string  `json:"f_name" validate:"required,min=2,max=100"`
	LastName    *string  `json:"l_name" validate:"required,min=2,max=100"`
	Email   *string  `json:"email" validation:"email,required"`
	Phone   *string  `json:"phone" validation:"required"`
	Address Address `json:"address" validation:"required"`
	Token   *string  `json:"token"`
	Refresh_token *string `json:"refresh_token"`
	Created_at  time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	User_id  string
}