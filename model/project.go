package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Project struct {
	Id    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Image string             `json:"image" bson:"image"`
	Url   string             `json:"url" bson:"url"`
}
