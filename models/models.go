package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title string             `json:"title" bson:"title,omitempty"`
	Body  string             `json:"body" bson:"body,omitempty"`
}
