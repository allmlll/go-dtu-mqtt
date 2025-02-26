package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Center struct {
	Id    primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Index int                `json:"index" bson:"index"`
}
