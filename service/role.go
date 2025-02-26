package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
	"ruitong-new-service/util"
)

type roleService struct {
}

var Role roleService

func (r *roleService) Add(role *model.Role) (any, error) {
	count, err := global.RoleColl.CountDocuments(context.TODO(), bson.M{"name": role.Name})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	if count != 0 {
		return nil, errors.New("角色名已存在")
	}
	_, err = global.RoleColl.InsertOne(context.TODO(), role)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	return nil, nil
}
func (r *roleService) Update(role *model.Role) (any, error) {
	var count int64
	var err error
	//判断更新的api是否存在
	for _, api := range role.Apis {
		count, err = global.ApiColl.CountDocuments(context.TODO(), api)
		if err != nil {
			global.Log.Error(err.Error())
			return nil, err
		}
		if count == 0 {
			return nil, errors.New("api不存在")
		}
	}

	filter := bson.M{"_id": role.Id}
	update := bson.M{"$set": role}
	global.RoleColl.FindOneAndUpdate(context.TODO(), filter, update)
	return nil, nil
}

func (r *roleService) Get() (any, error) {
	cur, err := global.RoleColl.Find(context.TODO(), bson.M{})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	var roles []model.Role
	err = cur.All(context.TODO(), &roles)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	return roles, nil
}

func (r *roleService) GetCascade() (any, error) {
	var cascades []model.Cascade

	opts := options.Find().SetSort(bson.M{"index": 1})
	cur, err := global.CenterColl.Find(context.TODO(), bson.M{}, opts)
	err = cur.All(context.TODO(), &cascades)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	for i, c := range cascades {
		cur, err = global.GroupColl.Find(context.TODO(), bson.M{"center": c.Value}, opts)
		var cascades1 []model.Cascade
		err = cur.All(context.TODO(), &cascades1)
		if err != nil {
			global.Log.Error(err.Error())
			return nil, err
		}
		for j, c2 := range cascades1 {
			cur, err = global.DeviceColl.Find(context.TODO(), bson.M{"group": c2.Value})
			var cascades2 []model.Cascade
			err = cur.All(context.TODO(), &cascades2)
			if err != nil {
				global.Log.Error(err.Error())
				return nil, err
			}
			cascades1[j].Children = cascades2
		}
		cascades[i].Children = cascades1
	}

	return cascades, nil
}

func (r *roleService) PaginationFind(page *util.Page, name string) (any, error) {
	opts := page.GetOpts()
	var filter bson.M
	if name != "" {
		regexString := fmt.Sprintf(".*%v.*", name)
		filter = bson.M{
			"name": primitive.Regex{Pattern: regexString},
		}
	}

	cur, err := global.RoleColl.Find(context.TODO(), filter, opts)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	var roles []model.Role
	err = cur.All(context.TODO(), &roles)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	count, err := global.RoleColl.CountDocuments(context.TODO(), filter)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	return gin.H{
		"roles": roles,
		"count": count,
	}, nil
}
