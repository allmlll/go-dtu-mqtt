package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
	"ruitong-new-service/util"

	"strings"
	"time"
)

type latestDataService struct {
}

var LatestData latestDataService

func (l *latestDataService) Create(latestData model.LatestData) error {
	count, err := global.LatestDataColl.CountDocuments(context.TODO(), bson.M{"code": latestData.Code})
	if count != 0 {
		return errors.New("设备ID已存在")
	}

	_, err = global.LatestDataColl.InsertOne(context.TODO(), latestData)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}

	return nil
}

func (l *latestDataService) BgiScreenMap() (any, error) {
	cur, err := global.LatestDataColl.Find(context.TODO(), bson.M{})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	var latestDatas []model.LatestData
	err = cur.All(context.TODO(), &latestDatas)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	data := make(map[string]model.LatestData, len(latestDatas))
	for _, latestData := range latestDatas {
		data[latestData.Code] = latestData
	}

	return data, nil
}

func (l *latestDataService) WebHome(page util.Page) (any, error) {
	cur, err := global.LatestDataColl.Find(context.TODO(), bson.M{}, page.GetOpts())
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	var latestDatas []model.LatestData
	err = cur.All(context.TODO(), &latestDatas)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	documents, err := global.LatestDataColl.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	return gin.H{
		"count": documents,
		"data":  latestDatas,
	}, nil
}

func (l *latestDataService) BigScreenCems(conn *websocket.Conn) error {
	// 获取cems小时表中每个code的最后一条数据
	var cemsHourDatas []model.DeviceData
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
	err = cur.All(context.TODO(), &cemsHourDatas)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}

	// 获取并整理小时数据
	// code-key-value
	hourDataMap := make(map[string]map[string]string)
	var codes []string
	for _, deviceData := range cemsHourDatas {
		// 整理设备code便于查找最新数据
		codes = append(codes, deviceData.Code)

		// 解析[]model.Data为map,去除多余数据,修改key值,便于查找
		keyMap := make(map[string]string)
		for _, keyData := range deviceData.Data {
			if strings.Contains(keyData.Key, "平均") {
				newKey := strings.Replace(keyData.Key, "平均值", "", -1)
				keyMap[newKey] = keyData.Value
			}
		}
		hourDataMap[deviceData.Code] = keyMap
	}

	// 发送历史最新数据
	cur, err = global.LatestDataColl.Find(context.TODO(), bson.M{"code": bson.M{"$in": codes}})
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}

	for cur.Next(context.TODO()) {
		var latestData model.LatestData
		err = cur.Decode(&latestData)
		if err != nil {
			global.Log.Error(err.Error())
			return err
		}

		keyMap := hourDataMap[latestData.Code]
		for _, data := range latestData.Data {
			msg := []string{latestData.Name, data.Key, data.Value, keyMap[data.Key]}
			marshal, _ := json.Marshal(msg)
			err = conn.WriteMessage(websocket.TextMessage, marshal)
			if err != nil {
				// 应该是连接关闭
				return nil
			}
		}
	}

	// 从redis获得最新数据发给前端
	var topics []string
	for _, code := range codes {
		topics = append(topics, global.LatestDataKeyPrefix+code)
	}

	subscribe := global.Redis.Subscribe(context.TODO(), topics...)
	defer subscribe.Close()
	for {
		message, err := util.RedisReceiveMessageTimeOut(subscribe, context.TODO(), time.Minute*5)
		if err != nil {
			global.Log.Error(err.Error())
			return err
		}

		var latestData model.LatestData
		_ = json.Unmarshal([]byte(message.Payload), &latestData)

		keyMap := hourDataMap[latestData.Code]
		for _, data := range latestData.Data {
			msg := []string{latestData.Name, data.Key, data.Value, keyMap[data.Key]}
			marshal, _ := json.Marshal(msg)
			err = conn.WriteMessage(websocket.TextMessage, marshal)
			if err != nil {
				// 应该是连接关闭
				return nil
			}
		}
	}
}

func (l *latestDataService) Delete(code string) error {
	_, err := global.LatestDataColl.DeleteOne(context.TODO(), bson.M{"code": code})
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}
	return nil
}

func (l *latestDataService) UpdateData(code string, datas []model.Data) error {
	filter := bson.M{"code": code}
	update := bson.M{
		"$set": bson.M{
			"data":       datas,
			"updateTime": time.Now().Format("2006-01-02 15:04:05"),
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var latestData model.LatestData
	err := global.LatestDataColl.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&latestData)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}

	// 跟新后的数据发送到redis
	marshal, _ := json.Marshal(latestData)
	global.Redis.Publish(context.TODO(), global.LatestDataKeyPrefix+latestData.Code, marshal)
	return nil
}

func (l *latestDataService) Update(code string, name string, timeInterval *[]int, checkTime int) error {
	filter := bson.M{"code": code}
	update := bson.M{"name": name, "timeInterval": timeInterval, "checkTime": checkTime}
	_, err := global.LatestDataColl.UpdateOne(context.TODO(), filter, bson.M{"$set": update})

	return err
}

func (l *latestDataService) CheckStatus() {
	// 查出所有数据
	var latestDatas []model.LatestData
	cur, err := global.LatestDataColl.Find(context.TODO(), bson.M{})
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	err = cur.All(context.TODO(), &latestDatas)
	if err != nil {
		global.Log.Error(err.Error())
		return
	}

	// 判断5分钟前有没有更新数据
	timeNow := time.Now()
	var offlineCodes, onlineCodes []string
	for _, latestData := range latestDatas {
		ago := timeNow.Add(time.Duration(-latestData.CheckTime) * time.Minute).Format("2006-01-02 15:04:05")
		if latestData.UpdateTime < ago {
			offlineCodes = append(offlineCodes, latestData.Code)
		} else {
			onlineCodes = append(onlineCodes, latestData.Code)
		}
	}
	if offlineCodes != nil {
		// 没有更新过数据的设备设为离线
		_, err = global.LatestDataColl.UpdateMany(context.TODO(), bson.M{
			"code": bson.M{
				"$in": offlineCodes,
			},
		}, bson.M{
			"$set": bson.M{"status": "离线"},
		})
		if err != nil {
			global.Log.Error(err.Error())
			return
		}
	}

	if onlineCodes != nil {
		// 更新过数据的设备设为在线
		_, err = global.LatestDataColl.UpdateMany(context.TODO(), bson.M{
			"code": bson.M{
				"$in": onlineCodes,
			},
		}, bson.M{
			"$set": bson.M{"status": "在线"},
		})
		if err != nil {
			global.Log.Error(err.Error())
			return
		}
	}

}
