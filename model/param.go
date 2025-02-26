package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Param struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Topic       string             `json:"topic" bson:"topic"`             // 发布的topic
	Payload     map[string]string  `json:"payload" bson:"payload"`         // 发布的内容
	Interval    int64              `json:"interval" bson:"interval"`       // 发布的间隔(s)
	LastPublish int64              `json:"lastPublish" bson:"lastPublish"` // 最后发布时间的时间戳
}
