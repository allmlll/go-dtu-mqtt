package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ruitong-new-service/global"
	"ruitong-new-service/model"

	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type parseData struct {
}

var ParseData parseData

type jsonType struct {
	DeviceCode string         `json:"did"`
	Data       map[string]any `json:"data"`
}

func (p *parseData) Json(s string) {
	defer func() {
		if err := recover(); err != nil {
			global.Log.Error(err.(error).Error() + "dtu解析错误 原始数据:" + s)
		}
	}()

	var jsonData jsonType
	if err := json.Unmarshal([]byte(s), &jsonData); err != nil {
		global.Log.Error(err.Error())
		return
	}

	// 生成data
	var datas []model.Data
	for k, v := range jsonData.Data {
		var value string
		switch reflect.TypeOf(v).Kind() {
		case reflect.Float64:
			value = strconv.FormatFloat(v.(float64), 'g', -1, 64)
		case reflect.String:
			value = v.(string)
		default:
			global.Log.Error("mqtt数据类型未知")
		}
		datas = append(datas, model.Data{
			Key:   k,
			Value: value,
		})
	}
	var device model.Device
	err := global.DeviceColl.FindOne(context.TODO(), bson.M{"code": jsonData.DeviceCode}).Decode(&device)
	if err != nil {
		global.Log.Error(err.Error())
	}
	if device.Sort == nil || len(device.Sort) == 0 {
		sort.Slice(datas, func(i, j int) bool {
			return datas[i].Key < datas[j].Key
		})
	} else {
		var datasNow []model.Data
		for _, key := range device.Sort {
			for i, data := range datas {
				if data.Key == key {
					datasNow = append(datasNow, data)
					datas = append(datas[:i], datas[i+1:]...)
					break
				}
			}
		}
		sort.Slice(datas, func(i, j int) bool {
			return datas[i].Key < datas[j].Key
		})
		datasNow = append(datasNow, datas...)
		datas = datasNow
	}

	now := time.Now()
	deviceData := model.DeviceData{
		Code:       jsonData.DeviceCode,
		Data:       datas,
		CreateTime: now.Format("2006-01-02 15:04") + ":00",
	}

	bytes, err := json.Marshal(deviceData)
	if err != nil {
		global.Log.Error(err.Error())
	}
	global.LogInfo.Info(string(bytes))
	p.insertDtuData(&deviceData)

	if err = LatestData.UpdateData(deviceData.Code, deviceData.Data); err != nil {
		global.Log.Error(err.Error())
	}
	p.updateDeviceSort(&deviceData, device.Sort)
}

func (p *parseData) Cems(s string, code string) {
	defer func() {
		if err := recover(); err != nil {
			global.Log.Error(err.(error).Error() + "cems解析错误 原始数据:" + s)
		}
	}()

	if s[0:2] != "##" {
		return
	}

	now := time.Now()
	var signal string
	var datas []model.Data

	parts := strings.Split(s, "&&")
	// 获取signal
	kvList := strings.Split(parts[0], ";")
	for _, kv := range kvList {
		kv := strings.Split(kv, "=")
		key := kv[0]
		if key == "CN" {
			signal = kv[1]
			break
		}
	}
	// 循环数据kv
	units := strings.Split(parts[1], ";")
	for _, unit := range units {
		kvList = strings.Split(unit, ",")
		for _, kv := range kvList {
			kv := strings.Split(kv, "=")
			if kv[0] == "DataTime" {
				if signal == "2061" || signal == "2031" {
					_, err := global.Redis.Get(context.TODO(), "repeatJudge:"+code+":"+signal+":"+kv[1]).Result()
					if errors.Is(err, redis.Nil) {
						if signal == "2061" {
							_, err = global.Redis.Set(context.TODO(), "repeatJudge:"+code+":"+signal+":"+kv[1], 0, 3*time.Hour).Result()
						} else {
							_, err = global.Redis.Set(context.TODO(), "repeatJudge:"+code+":"+signal+":"+kv[1], 0, 3*24*time.Hour).Result()
						}
						if err != nil {
							global.Log.Error(err.Error())
							return
						}
					} else if err != nil {
						global.Log.Error(err.Error())
						return
					} else {
						return
					}
				}
				continue
			}
			if strings.Contains(kv[0], "Flag") {
				continue
			}
			datas = append(datas, model.Data{
				Key:   kv[0],
				Value: kv[1],
			})
		}
	}

	var device model.Device
	err := global.DeviceColl.FindOne(context.TODO(), bson.M{"code": code}).Decode(&device)
	if err != nil {
		global.Log.Error(err.Error())
	}
	// 对datas排序
	if device.Sort == nil || len(device.Sort) == 0 {
		sort.Slice(datas, func(i, j int) bool {
			return datas[i].Key < datas[j].Key
		})
	} else {
		var datasNow []model.Data
		for _, key := range device.Sort {
			for i, data := range datas {
				if data.Key == key {
					datasNow = append(datasNow, data)
					datas = append(datas[:i], datas[i+1:]...)
					break
				}
			}
		}
		sort.Slice(datas, func(i, j int) bool {
			return datas[i].Key < datas[j].Key
		})
		datasNow = append(datasNow, datas...)
		datas = datasNow
	}

	// 从redis读取协议解析
	for i, data := range datas {
		value, err := global.Redis.Get(context.TODO(), global.CemsKeyPrefix+data.Key).Result()
		if err != nil {
			global.Log.Error(err.Error() + "key:" + data.Key)
			continue
		}
		datas[i].Key = value
	}

	// 组装数据
	deviceData := &model.DeviceData{
		Code:       code,
		Data:       datas,
		CreateTime: now.Format("2006-01-02 15:04") + ":00",
	}

	// 插入数据库
	p.insertCemsData(deviceData, signal)

	// 更新最新数据
	if signal == "2011" {
		if err := LatestData.UpdateData(deviceData.Code, deviceData.Data); err != nil {
			global.Log.Error(err.Error())
		}
	}
	if signal == "2051" {
		p.updateDeviceSort(deviceData, device.Sort)
	}
}

func (p *parseData) insertDtuData(deviceData *model.DeviceData) {
	coll, ok := global.CollMap[0].Load(deviceData.Code)
	if !ok {
		global.Log.Error(deviceData.Code + "数据库不存在")
		return
	}
	_, err := coll.(*mongo.Collection).InsertOne(context.Background(), deviceData)
	if err != nil {
		global.Log.Error(err.Error())
	}
}

func (p *parseData) insertCemsData(deviceData *model.DeviceData, signal string) {
	var coll *mongo.Collection
	switch signal {
	case "2011":
		coll = global.CemsRealtimeColl
	case "2051":
		coll = global.CemsMinuteColl
	case "2061":
		coll = global.CemsHourColl
	case "2031":
		coll = global.CemsDayColl
	default:
		global.Log.Error("cems 未知时间粒度表示:" + signal)
		return
	}
	_, err := coll.InsertOne(context.TODO(), deviceData)

	if err != nil {
		global.Log.Error(err.Error())
	}
	//更新大屏redis
	if signal == "2061" {
		p.insertRedis(deviceData.Code, coll)
	}

}

func (p *parseData) insertRedis(code string, coll *mongo.Collection) {
	//存入大屏redis

	key := "Screen:" + code
	cur, err := coll.Find(context.TODO(), bson.M{"code": code}, options.Find().SetSort(bson.M{"createTime": -1}).SetLimit(24))
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	var deviceDatas []model.DeviceData
	err = cur.All(context.TODO(), &deviceDatas)
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	marshal, err := json.Marshal(deviceDatas)
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	err = global.Redis.Set(context.TODO(), key, marshal, 0).Err()
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	return
}

func (p *parseData) updateDeviceSort(deviceData *model.DeviceData, sort []string) {
	if len(deviceData.Data) <= len(sort) {
		return
	}
	var sortString []string
	for _, data := range deviceData.Data {
		sortString = append(sortString, data.Key)
	}
	update := bson.M{"sort": sortString}
	if len(sort) == 0 {
		update["showKeys"] = sortString
	}
	_, err := global.DeviceColl.UpdateOne(context.TODO(), bson.M{"code": deviceData.Code}, bson.M{"$set": update})
	if err != nil {
		global.Log.Error(err.Error())
	}

	return
}
