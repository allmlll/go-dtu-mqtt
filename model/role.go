package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role struct {
	Id         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string               `json:"name" bson:"name"`
	Apis       []Api                `json:"apis" bson:"apis"`
	Centers    []primitive.ObjectID `json:"centers" bson:"centers"`
	Codes      []string             `json:"codes" bson:"codes"`
	FirstPage  string               `json:"firstPage" bson:"firstPage"`
	RoleRoutes string               `json:"roleRoutes" bson:"roleRoutes"`
	Desc       string               `json:"desc" bson:"desc"` //角色描述
}
