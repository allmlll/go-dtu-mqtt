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

type centerService struct {
}

var Center centerService

func (c *centerService) Add(center *model.Center) (any, error) {

	err := global.CenterColl.FindOne(context.TODO(), bson.M{"name": center.Name}).Decode(&model.Center{})
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("name重复")
	}

	count, err := global.CenterColl.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			center.Index = 0
		} else {
			global.Log.Error(err.Error())
			return nil, err
		}
	} else {
		center.Index = int(count)
	}

	_, err = global.CenterColl.InsertOne(context.TODO(), center)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	return nil, nil
}

func (c *centerService) PaginationFind(page *util.Page, name string) (any, error) {
	opts := page.GetOpts().SetSort(bson.M{"index": 1})
	var filter bson.M
	if name != "" {
		regexString := fmt.Sprintf(".*%v.*", name)
		filter = bson.M{
			"name": primitive.Regex{Pattern: regexString},
		}
	}

	cur, err := global.CenterColl.Find(context.TODO(), filter, opts)
	var centers []model.Center
	err = cur.All(context.TODO(), &centers)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	count, err := global.CenterColl.CountDocuments(context.TODO(), filter)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	return gin.H{
		"centers": centers,
		"count":   count,
	}, nil
}

func (c *centerService) Get() (any, error) {
	opts := options.Find().SetSort(bson.M{"index": 1})
	cur, err := global.CenterColl.Find(context.TODO(), bson.M{}, opts)
	var centers []model.Center
	err = cur.All(context.TODO(), &centers)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	centers = append(centers, model.Center{
		Id:   primitive.NilObjectID,
		Name: "默认中心",
	})
	return centers, nil
}

func (c *centerService) Update(center *model.Center) (any, error) {
	var DBCenter model.Center
	err := global.CenterColl.FindOne(context.TODO(), bson.M{"_id": center.Id}).Decode(&DBCenter)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	if center.Index > DBCenter.Index {
		_, err = global.GroupColl.UpdateMany(context.TODO(), bson.M{"index": bson.M{"$gt": DBCenter.Index}, "$lte": center.Index}, bson.M{"$inc": bson.M{"index": -1}})
		if err != nil {
			global.Log.Error(err.Error())
			return nil, err
		}
	} else if center.Index < DBCenter.Index {
		_, err = global.GroupColl.UpdateMany(context.TODO(), bson.M{"index": bson.M{"$gte": center.Index}, "$lt": DBCenter.Index}, bson.M{"$inc": bson.M{"index": 1}})
		if err != nil {
			global.Log.Error(err.Error())
			return nil, err
		}
	}

	_, err = global.CenterColl.UpdateOne(context.TODO(), bson.M{"_id": center.Id}, bson.M{"$set": center})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	return nil, nil
}

func (c *centerService) Delete(center *model.Center) (any, error) {
	filter := bson.M{
		"_id": center.Id,
	}
	_, err := global.CenterColl.DeleteOne(context.TODO(), filter)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	_, err = global.CenterColl.UpdateMany(context.TODO(), bson.M{"index": bson.M{"$gt": center.Index}}, bson.M{"$inc": bson.M{"index": -1}})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	_, err = global.GroupColl.UpdateMany(context.TODO(), bson.M{"center": center.Id}, bson.M{"$set": bson.M{"center": primitive.NilObjectID}})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	return nil, nil
}
