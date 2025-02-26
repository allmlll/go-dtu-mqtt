package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type BigScreen struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Position   string             `bson:"position" json:"position"` //左1 右1
	StateTypes []stateType        `bson:"stateTypes" json:"stateTypes"`
}

type stateType struct {
	Name  string `bson:"stateTypeName" json:"stateTypeName"`
	Items []Item `bson:"items" json:"items"`
}

type Item struct {
	NameList []string `bson:"nameList" json:"nameList"`
	Code     string   `bson:"code" json:"code"`
}
