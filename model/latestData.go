package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type LatestData struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Code         string             `json:"code" bson:"code"`
	Type         string             `json:"type" bson:"type"`
	Data         []Data             `json:"data" bson:"data"`
	Status       string             `json:"status" bson:"status"`
	TimeInterval []int              `bson:"timeInterval" json:"timeInterval"`
	CheckTime    int                `bson:"checkTime" json:"checkTime"`
	UpdateTime   string             `json:"updateTime" bson:"updateTime"`
}
