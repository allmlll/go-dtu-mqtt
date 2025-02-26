package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Device struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name         string             `bson:"name" json:"name"`
	Code         string             `bson:"code" json:"code"` //设备唯一标识
	Topic        string             `bson:"Topic" json:"Topic"`
	Type         string             `bson:"type" json:"type"` // dtu cems
	Group        primitive.ObjectID `bson:"group" json:"group"`
	Keys         []string           `bson:"keys" json:"keys"`
	Sort         []string           `bson:"sort" json:"sort"`
	ShowKeys     []string           `bson:"showKeys" json:"showKeys"`
	TimeInterval []int              `bson:"timeInterval" json:"timeInterval"`
	CheckTime    int                `bson:"checkTime" json:"checkTime"`
	CreateTime   string             `bson:"createTime" json:"createTime"`           //创建时间
	UpdateTime   string             `bson:"updateTime,omitempty" json:"updateTime"` //修改时间
}
