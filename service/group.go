package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
	"ruitong-new-service/util"
)

type groupService struct {
}

var Group groupService

func (g *groupService) Add(group *model.Group) (any, error) {

	err := global.GroupColl.FindOne(context.TODO(), bson.M{"name": group.Name}).Decode(&model.Group{})
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("name重复")
	}

	count, err := global.GroupColl.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			group.Index = 0
		} else {
			global.Log.Error(err.Error())
			return nil, err
		}
	} else {
		group.Index = int(count)
	}

	_, err = global.GroupColl.InsertOne(context.TODO(), group)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	return nil, nil
}

func (g *groupService) PaginationFind(page *util.Page, name, id string) (any, error) {
	opts := page.GetOpts().SetSort(bson.M{"index": 1})

	filter := bson.M{}
	if name != "" {
		regexString := fmt.Sprintf(".*%v.*", name)
		filter["name"] = primitive.Regex{Pattern: regexString}
	}
	if id != "" {
		oId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			global.Log.Error(err.Error())
			return nil, err
		}
		filter["center"] = oId
	}

	cur, err := global.GroupColl.Find(context.TODO(), filter, opts)
	var groups []model.Group
	err = cur.All(context.TODO(), &groups)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	var centers []model.Center
	cur, err = global.CenterColl.Find(context.TODO(), bson.M{})
	err = cur.All(context.TODO(), &centers)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	var resp []map[string]any
	for _, group := range groups {
		groupMap := make(map[string]any)
		groupMap["group"] = group
		if group.Center == primitive.NilObjectID {
			groupMap["center"] = "默认中心"
		} else {
			for _, center := range centers {
				if group.Center == center.Id {
					groupMap["center"] = center.Name
					break
				}
			}
		}
		resp = append(resp, groupMap)
	}

	count, err := global.GroupColl.CountDocuments(context.TODO(), filter)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	return gin.H{
		"groups": resp,
		"count":  count,
	}, nil
}

func (g *groupService) Get() (any, error) {
	opts := options.Find().SetSort(bson.M{"index": 1})
	cur, err := global.GroupColl.Find(context.TODO(), bson.M{}, opts)
	var groups []model.Group
	err = cur.All(context.TODO(), &groups)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	groups = append(groups, model.Group{
		Id:     primitive.NilObjectID,
		Name:   "默认分组",
		Center: primitive.NilObjectID,
	})
	return groups, nil
}

func (g *groupService) GetCascade() (any, error) {
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
		cascades[i].Children = cascades1
	}

	var cascades1 []model.Cascade
	cur, err = global.GroupColl.Find(context.TODO(), bson.M{"center": primitive.NilObjectID}, opts)
	err = cur.All(context.TODO(), &cascades1)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	cascades1 = append(cascades1, model.Cascade{
		Value: primitive.NilObjectID,
		Label: "默认分组",
	})

	cascades = append(cascades, model.Cascade{
		Value:    primitive.NilObjectID,
		Label:    "默认中心",
		Children: cascades1,
	})

	return cascades, nil
}

func (g *groupService) Update(group *model.Group) (any, error) {

	_, err := global.GroupColl.UpdateOne(context.TODO(), bson.M{"_id": group.Id}, bson.M{"$set": group})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	return nil, nil
}

func (g *groupService) Delete(group *model.Group) (any, error) {
	_, err := global.GroupColl.DeleteOne(context.TODO(), bson.M{"_id": group.Id})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	_, err = global.DeviceColl.UpdateMany(context.TODO(), bson.M{"group": group.Id}, bson.M{"$set": bson.M{"group": primitive.NilObjectID}})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	return nil, nil
}
