package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cascade struct {
	Value        primitive.ObjectID `json:"value" bson:"_id"`
	Label        string             `json:"label" bson:"name"`
	Code         string             `json:"code" bson:"code"`
	Type         string             `json:"type" bson:"type"`
	TimeInterval []int              `bson:"timeInterval" json:"timeInterval"`
	Children     []Cascade          `json:"children"`
}
