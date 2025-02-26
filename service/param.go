package service

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
	"ruitong-new-service/util"
	"time"
)

type paramService struct{}

var Param paramService

func (p *paramService) Add(param model.Param) (any, error) {
	param.LastPublish = time.Now().Unix()
	_, err := global.ParamColl.InsertOne(context.TODO(), param)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (p *paramService) Find(page *util.Page) (any, error) {
	count, err := global.ParamColl.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	var params []model.Param
	opts := page.GetOpts()
	cur, err := global.ParamColl.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	if err = cur.All(context.TODO(), &params); err != nil {
		return nil, err
	}

	return gin.H{
		"count":  count,
		"params": params,
	}, nil
}

func (p *paramService) Update(param model.Param) (any, error) {
	_, err := global.ParamColl.UpdateOne(context.TODO(), bson.M{"_id": param.Id}, bson.M{"$set": bson.M{
		"name":     param.Name,
		"topic":    param.Topic,
		"payload":  param.Payload,
		"interval": param.Interval,
	}})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (p *paramService) Delete(param model.Param) (any, error) {
	_, err := global.ParamColl.DeleteOne(context.TODO(), bson.M{"_id": param.Id})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (p *paramService) ParamPublish() {
	defer func() {
		if err := recover(); err != nil {
			err2 := err.(error)
			global.Log.Error(err2.Error())
		}
		go p.ParamPublish()
	}()

	for now := range time.Tick(time.Minute) {
		var params []model.Param
		cur, err := global.ParamColl.Find(context.TODO(), bson.M{})
		if err != nil {
			global.Log.Error(err.Error())
			continue
		}
		if err := cur.All(context.TODO(), &params); err != nil {
			global.Log.Error(err.Error())
			continue
		}

		for _, param := range params {
			if now.Unix()-param.LastPublish >= param.Interval {
				// 发布参数
				marshal, err := json.Marshal(param.Payload)
				if err != nil {
					global.Log.Error(err.Error())
					continue
				}
				global.MqttClient.Publish(param.Topic, 0, false, marshal)

				// 修改最后发布时间
				_, err = global.ParamColl.UpdateOne(context.TODO(), bson.M{"_id": param.Id}, bson.M{"$set": bson.M{
					"lastPublish": now.Unix(),
				}})
				if err != nil {
					global.Log.Error(err.Error())
					continue
				}
			}
		}
	}
}
