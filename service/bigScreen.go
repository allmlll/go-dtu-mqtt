package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
	"ruitong-new-service/util"
	"strconv"
	"strings"
	"time"
)

type res struct {
	Codes []string `json:"codes"`
}
type bigScreenService struct {
	Res             res
	GraphiteWSCodes []string
}

var BigScreen = bigScreenService{
	GraphiteWSCodes: []string{},
}

type bigScreenData struct {
	Details []detail `json:"detail"`
}
type detail struct {
	Code       string             `json:"code"`
	DeviceData []model.DeviceData `json:"deviceData"`
}
type Table struct {
	Key     string `bson:"key" json:"key"`
	Now     string `bson:"now" json:"now"`
	Five    string `bson:"five" json:"five"`
	Fifteen string `bson:"fifteen" json:"fifteen"`
}

type Table1 struct {
	Key     string `bson:"key" json:"key"`
	Now     string `bson:"now" json:"now"`
	Fifteen string `bson:"fifteen" json:"fifteen"`
	Hour    string `bson:"hour" json:"hour"`
}

func (*bigScreenService) Get(r res) (any, error) {
	var bigs []detail
	for _, code := range r.Codes {
		key := "Screen:" + code
		result, err := global.Redis.Get(context.Background(), key).Result()
		if err != nil {
			global.Log.Error(err.Error())
			continue
		}

		var deviceDatas []model.DeviceData
		err = json.Unmarshal([]byte(result), &deviceDatas)
		if err != nil {
			global.Log.Error(err.Error())
			continue
		}
		bigs = append(bigs, detail{
			Code:       code,
			DeviceData: deviceDatas,
		})
	}
	return bigScreenData{Details: bigs}, nil
}
func (*bigScreenService) GetAll() (any, error) {
	var bigs []detail
	keys, err := global.Redis.Keys(context.TODO(), "Screen:*").Result()
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	for _, key := range keys {
		result, err := global.Redis.Get(context.Background(), key).Result()
		if err != nil {
			global.Log.Error(err.Error())
			continue
		}

		var deviceDatas []model.DeviceData
		err = json.Unmarshal([]byte(result), &deviceDatas)
		if err != nil {
			global.Log.Error(err.Error())
			continue
		}
		bigs = append(bigs, detail{
			Code:       key[7:],
			DeviceData: deviceDatas,
		})

	}
	return bigScreenData{Details: bigs}, nil
}

func (*bigScreenService) GetCount() (any, error) {
	dtuCount, err := global.Redis.Get(context.TODO(), "count:dtuCount").Result()
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	cemsCount, err := global.Redis.Get(context.TODO(), "count:cemsCount").Result()
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	return gin.H{
		"cemsCount": dtuCount,
		"dtuCount":  cemsCount,
	}, nil
}

func (*bigScreenService) GetGraphiteWSData(conn *websocket.Conn) error {
	var latestDatas []model.LatestData
	cur, err := global.LatestDataColl.Find(context.TODO(), bson.M{"code": bson.M{"$in": BigScreen.GraphiteWSCodes}})
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}
	err = cur.All(context.TODO(), &latestDatas)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}

	latestDataMap := make(map[string]model.LatestData)
	for _, data := range latestDatas {
		latestDataMap[data.Code] = data
	}

	result := make(map[string]any)
	result["graConfiguration"] = latestDataMap[BigScreen.GraphiteWSCodes[0]].Data

	if latestDataMap[BigScreen.GraphiteWSCodes[0]].Status == "在线" {
		result["graDeviceStatus"] = true
	} else {
		result["graDeviceStatus"] = false
	}
	tableMap := make(map[string]*Table, 0)

	for _, data := range latestDataMap[BigScreen.GraphiteWSCodes[1]].Data {
		tableMap[data.Key] = &Table{
			Key:  data.Key,
			Five: data.Value,
		}
	}
	for _, data := range latestDataMap[BigScreen.GraphiteWSCodes[2]].Data {
		if tableMap[data.Key] == nil {
			tableMap[data.Key] = &Table{
				Key:     data.Key,
				Fifteen: data.Value,
			}
		} else {
			tableMap[data.Key].Fifteen = data.Value
		}
	}
	table := make([]Table, 0)
	for _, data := range latestDataMap[BigScreen.GraphiteWSCodes[0]].Data {
		if tableMap[data.Key] == nil {
			tableMap[data.Key] = &Table{
				Key: data.Key,
				Now: data.Value,
			}
		} else {
			tableMap[data.Key].Now = data.Value
		}
		table = append(table, *tableMap[data.Key])
	}
	result["table"] = table

	err = conn.WriteJSON(result)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}
	subscribe := global.Redis.Subscribe(context.TODO(), global.LatestDataKeyPrefix+BigScreen.GraphiteWSCodes[0])
	defer subscribe.Close()
	for {
		message, err := util.RedisReceiveMessageTimeOut(subscribe, context.TODO(), time.Minute*5)
		if err != nil {
			global.Log.Error(err.Error())
			return err
		}

		var latestData model.LatestData
		_ = json.Unmarshal([]byte(message.Payload), &latestData)
		result["graConfiguration"] = latestData.Data
		table = make([]Table, 0)
		for _, data := range latestData.Data {
			tableMap[data.Key].Now = data.Value
			table = append(table, *tableMap[data.Key])
		}
		err = conn.WriteJSON(result)
		if err != nil {
			global.Log.Error(err.Error())
			return err
		}
	}
}

func (*bigScreenService) GetTunnelKilnWSData(conn *websocket.Conn) error {
	code := "211124071146"
	var latestData model.LatestData
	err := global.LatestDataColl.FindOne(context.TODO(), bson.M{"code": code}).Decode(&latestData)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}
	fifteenColl, _ := global.CollMap[1].Load(code)
	opt := options.FindOne().SetSort(bson.M{"createTime": -1})
	var fifteenDeviceData model.DeviceData
	err = fifteenColl.(*mongo.Collection).FindOne(context.TODO(), bson.M{}, opt).Decode(&fifteenDeviceData)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}

	hourColl, _ := global.CollMap[2].Load(code)
	var hourDeviceData model.DeviceData
	err = hourColl.(*mongo.Collection).FindOne(context.TODO(), bson.M{}, opt).Decode(&hourDeviceData)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}

	result := make(map[string]any)
	result["graConfiguration"] = latestData.Data

	if latestData.Status == "在线" {
		result["graDeviceStatus"] = true
	} else {
		result["graDeviceStatus"] = false
	}
	tableMap := make(map[string]*Table1, 0)

	for _, data := range hourDeviceData.Data {
		tableMap[data.Key] = &Table1{
			Key:  data.Key,
			Hour: data.Value,
		}
	}
	for _, data := range fifteenDeviceData.Data {
		if tableMap[data.Key] == nil {
			tableMap[data.Key] = &Table1{
				Key:     data.Key,
				Fifteen: data.Value,
			}
		} else {
			tableMap[data.Key].Fifteen = data.Value
		}
	}
	table := make([]Table1, 0)
	for _, data := range latestData.Data {
		if tableMap[data.Key] == nil {
			tableMap[data.Key] = &Table1{
				Key: data.Key,
				Now: data.Value,
			}
		} else {
			tableMap[data.Key].Now = data.Value
		}
		table = append(table, *tableMap[data.Key])
	}
	result["table"] = table

	err = conn.WriteJSON(result)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}
	subscribe := global.Redis.Subscribe(context.TODO(), global.LatestDataKeyPrefix+code)
	defer subscribe.Close()
	for {
		message, err := util.RedisReceiveMessageTimeOut(subscribe, context.TODO(), time.Minute*5)
		if err != nil {
			global.Log.Error(err.Error())
			return err
		}

		var latestData model.LatestData
		_ = json.Unmarshal([]byte(message.Payload), &latestData)
		result["graConfiguration"] = latestData.Data
		table = make([]Table1, 0)
		for _, data := range latestData.Data {
			tableMap[data.Key].Now = data.Value
			table = append(table, *tableMap[data.Key])
		}
		err = conn.WriteJSON(result)
		if err != nil {
			global.Log.Error(err.Error())
			return err
		}
	}
}

func (*bigScreenService) GetLargePressesWSData(conn *websocket.Conn) error {
	code := ""
	var latestData model.LatestData
	err := global.LatestDataColl.FindOne(context.TODO(), bson.M{"code": code}).Decode(&latestData)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}
	fifteenColl, _ := global.CollMap[1].Load(code)
	opt := options.FindOne().SetSort(bson.M{"createTime": -1})
	var fifteenDeviceData model.DeviceData
	err = fifteenColl.(*mongo.Collection).FindOne(context.TODO(), bson.M{}, opt).Decode(&fifteenDeviceData)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}

	hourColl, _ := global.CollMap[2].Load(code)
	var hourDeviceData model.DeviceData
	err = hourColl.(*mongo.Collection).FindOne(context.TODO(), bson.M{}, opt).Decode(&hourDeviceData)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}

	result := make(map[string]any)
	result["graConfiguration"] = latestData.Data

	if latestData.Status == "在线" {
		result["graDeviceStatus"] = true
	} else {
		result["graDeviceStatus"] = false
	}
	tableMap := make(map[string]*Table1, 0)

	for _, data := range hourDeviceData.Data {
		tableMap[data.Key] = &Table1{
			Key:  data.Key,
			Hour: data.Value,
		}
	}
	for _, data := range fifteenDeviceData.Data {
		if tableMap[data.Key] == nil {
			tableMap[data.Key] = &Table1{
				Key:     data.Key,
				Fifteen: data.Value,
			}
		} else {
			tableMap[data.Key].Fifteen = data.Value
		}
	}
	table := make([]Table1, 0)
	for _, data := range latestData.Data {
		if tableMap[data.Key] == nil {
			tableMap[data.Key] = &Table1{
				Key: data.Key,
				Now: data.Value,
			}
		} else {
			tableMap[data.Key].Now = data.Value
		}
		table = append(table, *tableMap[data.Key])
	}
	result["table"] = table

	err = conn.WriteJSON(result)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}
	subscribe := global.Redis.Subscribe(context.TODO(), global.LatestDataKeyPrefix+code)
	defer subscribe.Close()
	for {
		message, err := util.RedisReceiveMessageTimeOut(subscribe, context.TODO(), time.Minute*5)
		if err != nil {
			global.Log.Error(err.Error())
			return err
		}

		var latestData model.LatestData
		_ = json.Unmarshal([]byte(message.Payload), &latestData)
		result["graConfiguration"] = latestData.Data
		table = make([]Table1, 0)
		for _, data := range latestData.Data {
			tableMap[data.Key].Now = data.Value
			table = append(table, *tableMap[data.Key])
		}
		err = conn.WriteJSON(result)
		if err != nil {
			global.Log.Error(err.Error())
			return err
		}
	}
}

func (*bigScreenService) GetDippingWSData(conn *websocket.Conn) error {
	code := "211124071145"
	var latestData model.LatestData
	err := global.LatestDataColl.FindOne(context.TODO(), bson.M{"code": code}).Decode(&latestData)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}
	fifteenColl, _ := global.CollMap[1].Load(code)
	opt := options.FindOne().SetSort(bson.M{"createTime": -1})
	var fifteenDeviceData model.DeviceData
	err = fifteenColl.(*mongo.Collection).FindOne(context.TODO(), bson.M{}, opt).Decode(&fifteenDeviceData)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}

	hourColl, _ := global.CollMap[2].Load(code)
	var hourDeviceData model.DeviceData
	err = hourColl.(*mongo.Collection).FindOne(context.TODO(), bson.M{}, opt).Decode(&hourDeviceData)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}

	result := make(map[string]any)
	result["graConfiguration"] = latestData.Data

	if latestData.Status == "在线" {
		result["graDeviceStatus"] = true
	} else {
		result["graDeviceStatus"] = false
	}
	tableMap := make(map[string]*Table1, 0)

	for _, data := range hourDeviceData.Data {
		tableMap[data.Key] = &Table1{
			Key:  data.Key,
			Hour: data.Value,
		}
	}
	for _, data := range fifteenDeviceData.Data {
		if tableMap[data.Key] == nil {
			tableMap[data.Key] = &Table1{
				Key:     data.Key,
				Fifteen: data.Value,
			}
		} else {
			tableMap[data.Key].Fifteen = data.Value
		}
	}
	table := make([]Table1, 0)
	for _, data := range latestData.Data {
		if tableMap[data.Key] == nil {
			tableMap[data.Key] = &Table1{
				Key: data.Key,
				Now: data.Value,
			}
		} else {
			tableMap[data.Key].Now = data.Value
		}
		table = append(table, *tableMap[data.Key])
	}
	result["table"] = table

	err = conn.WriteJSON(result)
	if err != nil {
		global.Log.Error(err.Error())
		return err
	}
	subscribe := global.Redis.Subscribe(context.TODO(), global.LatestDataKeyPrefix+code)
	defer subscribe.Close()
	for {
		message, err := util.RedisReceiveMessageTimeOut(subscribe, context.TODO(), time.Minute*5)
		if err != nil {
			global.Log.Error(err.Error())
			return err
		}

		var latestData model.LatestData
		_ = json.Unmarshal([]byte(message.Payload), &latestData)
		result["graConfiguration"] = latestData.Data
		table = make([]Table1, 0)
		for _, data := range latestData.Data {
			tableMap[data.Key].Now = data.Value
			table = append(table, *tableMap[data.Key])
		}
		err = conn.WriteJSON(result)
		if err != nil {
			global.Log.Error(err.Error())
			return err
		}
	}
}

func (*bigScreenService) GetReportData(start, end, code string, interval int) (any, error) {
	if start == "" || end == "" {
		now := time.Now()
		end = now.Format("2006-01-02 15:04:05")
		start = now.Add(time.Hour * 24 * -1).Format("2006-01-02 15:04:05")
	}
	coll, ok := global.CollMap[interval].Load(code)
	if !ok {
		return nil, errors.New("code不存在")
	}
	filter := bson.M{"createTime": bson.M{
		"$gte": start,
		"$lte": end,
	}}
	opt := options.Find().SetSort(bson.M{"createTime": 1})
	cur, err := coll.(*mongo.Collection).Find(context.TODO(), filter, opt)
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
	dataMap := make(map[string][]string)
	for _, data := range datas {
		if dataMap["createTime"] == nil {
			dataMap["createTime"] = []string{}
		}
		dataMap["createTime"] = append(dataMap["createTime"], data.CreateTime)
		for _, datum := range data.Data {
			if dataMap[datum.Key] == nil {
				dataMap[datum.Key] = []string{}
			}
			dataMap[datum.Key] = append(dataMap[datum.Key], datum.Value)
		}
	}
	return dataMap, nil
}

func (*bigScreenService) GetHisStove() (any, error) {
	coll, _ := global.CollMap[0].Load(BigScreen.GraphiteWSCodes[3])
	var datas []model.DeviceData
	cur, err := coll.(*mongo.Collection).Find(context.TODO(), bson.M{}, options.Find().SetSort(bson.M{"createTime": -1}).SetLimit(5))
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	err = cur.All(context.TODO(), &datas)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	for i := range datas {
		datas[i].Code = strings.Split(datas[i].CreateTime, ":")[0]
		datas[i].Code = strings.ReplaceAll(datas[i].Code, "-", "")
		datas[i].Code = "SMH" + strings.ReplaceAll(datas[i].Code, " ", "")
	}

	return datas, nil
}

func (*bigScreenService) GetProfilingIngredientsData() (any, error) {
	code := "211124071140MR"
	coll, _ := global.CollMap[0].Load(code)
	now := time.Now()
	filter := bson.M{}
	dataMap := make(map[string][]string)
	dataMap["createTime"] = []string{}
	dataMap["锅数"] = []string{}
	dataMap["锅数"] = []string{}
	dataMap["总重"] = []string{}

	filter["createTime"] = bson.M{
		"$gte": now.Format("2006-01-02") + " 00:00:00",
	}
	var deviceDatas []model.DeviceData
	weight := 0.0
	for i := 0; i < 7; i++ {
		cur, err := coll.(*mongo.Collection).Find(context.TODO(), filter)
		if err != nil {
			global.Log.Error(err.Error())
			return nil, err
		}
		err = cur.All(context.TODO(), &deviceDatas)
		if err != nil {
			global.Log.Error(err.Error())
			return nil, err
		}

		dataMap["createTime"] = append(dataMap["createTime"], now.Format("2006-01-02")+" 00:00:00")
		dataMap["锅数"] = append(dataMap["锅数"], strconv.Itoa(len(deviceDatas)))
		for _, data := range deviceDatas {
			for _, datum := range data.Data {
				if datum.Key == "总重" {
					float, err := strconv.ParseFloat(datum.Value, 64)
					if err != nil {
						global.Log.Error(err.Error())
						return nil, err
					}
					weight += float
				}
			}
		}

		dataMap["总重"] = append(dataMap["总重"], strconv.FormatFloat(weight, 'g', -1, 64))

		filter["createTime"] = bson.M{
			"$gte": now.Add(time.Hour*-24).Format("2006-01-02") + " 00:00:00",
			"$lte": now.Format("2006-01-02") + " 00:00:00",
		}
		now = now.Add(time.Hour * -24)
	}
	return dataMap, nil
}
