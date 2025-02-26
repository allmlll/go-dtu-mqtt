package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
	"ruitong-new-service/util"
)

type cemsService struct {
}

var Cems cemsService

func (c *cemsService) Get(page *util.Page, interval int, start, end, code string) (any, error) {
	var coll *mongo.Collection
	switch interval {
	case 0:
		coll = global.CemsRealtimeColl
	case 1:
		coll = global.CemsMinuteColl
	case 2:
		coll = global.CemsHourColl
	case 3:
		coll = global.CemsDayColl
	}

	filter := bson.M{"code": code}
	if !(start == "" || end == "") {
		filter["createTime"] = bson.M{
			"$gte": start,
			"$lte": end,
		}
	}

	opt := page.GetOpts().SetSort(bson.M{"_id": -1})

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
	if interval != 0 && len(device.ShowKeys) != 0 && len(datas) != 0 {
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

func (c *cemsService) BigScreen(interval int, start, end, code string) (any, error) {
	var coll *mongo.Collection
	switch interval {
	case 0:
		coll = global.CemsRealtimeColl
	case 1:
		coll = global.CemsMinuteColl
	case 2:
		coll = global.CemsHourColl
	case 3:
		coll = global.CemsDayColl
	}

	filter := bson.M{"code": code}
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

func (c *cemsService) EveryCemsLatestData() (any, error) {
	// 获取cems小时表中每个code的最后一条数据
	var datas []model.DeviceData
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":          "$code",
				"lastDocument": bson.M{"$last": "$$ROOT"},
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$lastDocument",
			},
		},
	}
	cur, err := global.CemsHourColl.Aggregate(context.TODO(), pipeline)
	err = cur.All(context.TODO(), &datas)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	return gin.H{
		"datas": datas,
	}, nil
}
