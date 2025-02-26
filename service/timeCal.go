package service

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ruitong-new-service/global"
	"ruitong-new-service/model"

	"strconv"
	"time"
)

type timeTask struct {
	report []string
}

var TimeTask = timeTask{
	report: []string{},
}

func (t *timeTask) FifMin() {
	var insertData model.DeviceData
	global.CollMap[1].Range(func(code, coll any) bool {
		findColl, _ := global.CollMap[0].Load(code)
		err := findColl.(*mongo.Collection).FindOne(context.TODO(), bson.M{}, options.FindOne().SetSort(bson.M{"createTime": -1})).Decode(&insertData)
		if err != nil {
			global.Log.Error(err.Error())
			return true
		}
		startTime, err := time.Parse("2006-01-02 15:04:05", insertData.CreateTime)
		if err != nil {
			global.Log.Error(err.Error())
			return true
		}
		//超时则不计算
		if time.Now().Sub(startTime) > 5*time.Minute {
			return true
		}
		insertData.Id = primitive.ObjectID{}
		insertData.CreateTime = time.Now().Format("2006-01-02 15:04") + ":00"
		_, err = coll.(*mongo.Collection).InsertOne(context.TODO(), insertData)
		if err != nil {
			global.Log.Error(err.Error())
			return true
		}
		return true
	})
}

func (t *timeTask) Hour() {
	var insertData model.DeviceData

	global.CollMap[2].Range(func(code, coll any) bool {
		findColl, _ := global.CollMap[0].Load(code)
		err := findColl.(*mongo.Collection).FindOne(context.TODO(), bson.M{}, options.FindOne().SetSort(bson.M{"createTime": -1})).Decode(&insertData)
		if err != nil {
			global.Log.Error(err.Error())
			return true
		}
		//最近一条数据的时间
		startTime, err := time.Parse("2006-01-02 15:04:05", insertData.CreateTime)
		if err != nil {
			global.Log.Error(err.Error())
			return true
		}
		//超时则不计算
		if time.Now().Sub(startTime) > 30*time.Minute {
			return true
		}
		insertData.Id = primitive.ObjectID{}
		insertData.CreateTime = time.Now().Format("2006-01-02 15") + ":00:00"

		_, err = coll.(*mongo.Collection).InsertOne(context.TODO(), insertData)
		if err != nil {
			global.Log.Error(err.Error())
			return true
		}
		return true
	})
}

func (t *timeTask) HourAvg() {
	//初始时间
	now := time.Now().Format("2006-01-02 15") + ":00:00"
	oneHourAgo := time.Now().Add(-1*time.Hour).Format("2006-01-02 15") + ":00:00"

	var devices []model.Device
	cur, err := global.DeviceColl.Find(context.TODO(), bson.M{"type": "dtu"})
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	err = cur.All(context.TODO(), &devices)
	if err != nil {
		global.Log.Error(err.Error())
		return
	}

	var all []model.DeviceData
loop:
	for _, device := range devices {
		for _, timeInterval := range device.TimeInterval {
			if timeInterval == 3 {
				coll, ok := global.CollMap[0].Load(device.Code)
				if !ok {
					global.Log.Error(device.Code + "设备不存在")
					continue
				}
				cur, err = coll.(*mongo.Collection).Find(context.TODO(), bson.M{"createTime": bson.M{"$gte": oneHourAgo, "$lte": now}}, options.Find().SetSort(bson.M{"_id": -1}))
				if err != nil {
					global.Log.Error(err.Error())
					continue
				}
				err = cur.All(context.TODO(), &all)
				if err != nil {
					global.Log.Error(err.Error())
					continue
				}
				if len(all) == 0 {
					continue
				}
				//超时判断
				startTime, err := time.Parse("2006-01-02 15:04:05", all[len(all)-1].CreateTime)
				if err != nil {
					global.Log.Error(err.Error())
					continue loop
				}
				if time.Now().Sub(startTime) > 30*time.Minute {
					continue loop
				}

				keys := make(map[string]string)
				for _, key := range device.Keys {
					keys[key] = ""
				}

				//计算平均
				avg := all[0]
				for i, data := range all {
					if i == 0 {
						continue
					}
					for j, datum := range data.Data {
						if _, ok = keys[datum.Key]; !ok {
							continue
						}
						float1, err := strconv.ParseFloat(datum.Value, 64)
						float2, err := strconv.ParseFloat(avg.Data[j].Value, 64)
						if err != nil {
							global.Log.Error(err.Error())
							continue loop
						}
						avg.Data[j].Value = strconv.FormatFloat(float1+float2, 'g', -1, 64)
					}
				}
				for i, datum := range avg.Data {
					if _, ok = keys[datum.Key]; !ok {
						continue
					}
					float, err := strconv.ParseFloat(datum.Value, 64)
					if err != nil {
						global.Log.Error(err.Error())
						continue loop
					}
					datum.Value = strconv.FormatFloat(float/float64(len(all)), 'g', -1, 64)
					avg.Data[i] = datum
				}
				avg.Id = primitive.ObjectID{}
				avg.CreateTime = now

				insertColl1, _ := global.CollMap[3].Load(device.Code)

				insertColl := insertColl1.(*mongo.Collection)
				_, err = insertColl.InsertOne(context.TODO(), avg)
				if err != nil {
					global.Log.Error(err.Error())
					continue loop
				}
				t.bigScreenUpdate(insertColl)
			}
		}
	}
}

func (t *timeTask) bigScreenUpdate(coll *mongo.Collection) {
	option := options.Find()
	option.SetSort(bson.M{"createTime": -1})
	option.SetLimit(24)
	cur, err := coll.Find(context.TODO(), bson.M{}, option)

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
	err = global.Redis.Set(context.TODO(), "Screen:"+coll.Name(), marshal, 0).Err()
	if err != nil {
		global.Log.Error(err.Error())
		return
	}

}

func (t *timeTask) Day() {
	now := time.Now()
	cur, err := global.DeviceColl.Find(context.TODO(), bson.M{"code": bson.M{"$in": t.report}})
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	var devices []model.Device
	err = cur.All(context.TODO(), &devices)
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	opts := options.FindOne().SetSort(bson.M{"createTime": -1})
	filter := bson.M{
		"createTime": bson.M{
			"$gte": now.Format("2006-01-02 15") + ":00:00",
		},
		"data": bson.M{
			"key":   "日报表指示",
			"value": "1",
		},
	}
	var deviceData model.DeviceData
	for _, device := range devices {
		coll, ok := global.CollMap[0].Load(device.Code)
		if !ok {
			global.Log.Error(device.Code + "设备不存在")
			continue
		}
		err = coll.(*mongo.Collection).FindOne(context.TODO(), filter, opts).Decode(&deviceData)
		if err != nil {
			global.Log.Error(err.Error())
			continue
		}
		deviceData.CreateTime = now.Format("2006-01-02 15") + ":00:00"
		deviceData.Id = primitive.ObjectID{}
		_, err = global.ReportColl.InsertOne(context.TODO(), deviceData)
		if err != nil {
			global.Log.Error(err.Error())
			continue
		}
	}
}

func (t *timeTask) FifMinCulCount() {
	var center model.Center
	err := global.CenterColl.FindOne(context.TODO(), bson.M{"name": "环保数据中心"}).Decode(&center)
	if err != nil {
		global.Log.Error(err.Error())
		return
	}

	var groups []model.Group
	cur, err := global.GroupColl.Find(context.TODO(), bson.M{"center": center.Id})
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	err = cur.All(context.TODO(), &groups)
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	var ids []primitive.ObjectID
	var devices []model.Device
	for _, group := range groups {
		ids = append(ids, group.Id)
	}
	cur, err = global.DeviceColl.Find(context.TODO(), bson.M{"group": bson.M{"$in": ids}})
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	err = cur.All(context.TODO(), &devices)
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	var codes []string
	var dtuCount int64
	filter := bson.M{
		"createTime": bson.M{
			"$gte": time.Now().Add(-1*30*24*time.Hour).Format("2006-01-02 15:04") + ":00",
			"$lte": time.Now().Format("2006-01-02 15:04") + ":00",
		},
	}
	for _, device := range devices {
		if device.Type == "dtu" {
			coll, ok := global.CollMap[0].Load(device.Code)
			if !ok {
				continue
			}
			count, err := coll.(*mongo.Collection).CountDocuments(context.TODO(), filter)
			if err != nil {
				global.Log.Error(err.Error())
				return
			}
			dtuCount += count
		} else {
			codes = append(codes, device.Code)
		}
	}
	filter["code"] = bson.M{"$in": codes}
	cemsCount, err := global.CemsRealtimeColl.CountDocuments(context.TODO(), filter)
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	err = global.Redis.Set(context.TODO(), "count:dtuCount", dtuCount, 0).Err()
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
	err = global.Redis.Set(context.TODO(), "count:cemsCount", cemsCount, 0).Err()
	if err != nil {
		global.Log.Error(err.Error())
		return
	}
}
