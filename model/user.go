package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` // id
	Account  string             `json:"account" bson:"account"`            // 账号
	Name     string             `json:"name" bson:"name"`                  //姓名
	Phone    string             `json:"phone" bson:"phone"`
	Sex      string             `json:"sex" bson:"sex"`                       // 性别
	Password string             `json:"password" bson:"password"`             // 密码
	Avatar   string             `json:"avatar" bson:"avatar"`                 // 头像
	Role     primitive.ObjectID `json:"role,omitempty" bson:"role,omitempty"` // 角色外键
}
