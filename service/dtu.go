package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
	"ruitong-new-service/util"
)

type dtuService struct {
}

var Dtu dtuService

func (d *dtuService) Get(page *util.Page, interval int, start, end, code string) (any, error) {
	v, ok := global.CollMap[interval].Load(code)
	if !ok {
		return nil, errors.New("code不存在")
	}
	coll := v.(*mongo.Collection)

	filter := bson.M{}
	if !(start == "" || end == "") {
		filter["createTime"] = bson.M{
			"$gte": start,
			"$lte": end,
		}
	}
	opt := page.GetOpts().SetSort(bson.M{"_id": -1})
	//opt := options.Find().SetSkip(int64((page - 1) * size)).SetLimit(int64(size)).SetSort(bson.M{"_id": -1})

	cur, err := coll.Find(context.TODO(), filter, opt)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	var datas []model.DeviceData
	err = cur.All(context.TODO(), &datas)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	var device model.Device
	err = global.DeviceColl.FindOne(context.TODO(), bson.M{"code": code}).Decode(&device)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	if len(device.ShowKeys) != 0 && len(datas) != 0 {
		dataIndexMap := make(map[string]int)
		for i, data := range datas[0].Data {
			dataIndexMap[data.Key] = i
		}
		var dataNew []model.Data
		for _, key := range device.ShowKeys {
			dataNew = append(dataNew, datas[0].Data[dataIndexMap[key]])
		}
		datas[0].Data = dataNew
	}

	return gin.H{
		"data":  datas,
		"count": count,
	}, nil
}

func (d *dtuService) BigScreen(interval int, start, end, code string) (any, error) {
	v, ok := global.CollMap[interval].Load(code)
	if !ok {
		return nil, errors.New("code不存在")
	}
	coll := v.(*mongo.Collection)

	filter := bson.M{}
	if !(start == "" || end == "") {
		filter["createTime"] = bson.M{
			"$gte": start,
			"$lte": end,
		}
	}
	opt := options.Find().SetSort(bson.M{"_id": -1})

	cur, err := coll.Find(context.TODO(), filter, opt)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	var datas []model.DeviceData
	err = cur.All(context.TODO(), &datas)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	return gin.H{
		"data":  datas,
		"count": count,
	}, nil
}
