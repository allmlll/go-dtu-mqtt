package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Group struct {
	Id     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name   string             `json:"name" bson:"name"`
	Index  int                `json:"index" bson:"index"`
	Center primitive.ObjectID `json:"center" bson:"center"`
}
