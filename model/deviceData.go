package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeviceData struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Code       string             `bson:"code" json:"code"`
	Data       []Data             `bson:"data" json:"data"`                                 //详细数据信息
	CreateTime string             `bson:"createTime,omitempty" json:"createTime,omitempty"` //创建时间
}

type Data struct {
	Key   string `bson:"key" json:"key"`     //传感器能检测到的key值
	Value string `bson:"value" json:"value"` //value值
}
