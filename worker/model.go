package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Age          int                `json:"age" bson:"age"`
	MobileNumber string             `json:"mobile_number" bson:"mobileNumber"`
	Email        string             `json:"email" bson:"email"`
}
