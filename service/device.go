package service

import (
	"context"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
	"ruitong-new-service/util"
	"sort"

	"sync"
)

type deviceService struct {
	Topics sync.Map
}

var Device deviceService

func (d *deviceService) TopicInit() {
	var devices []model.Device
	cur, err := global.DeviceColl.Find(context.TODO(), bson.M{})
	if err != nil {
		global.Log.Error(err.Error())
	}
	if err = cur.All(context.TODO(), &devices); err != nil {
		global.Log.Error(err.Error())
	}
	for _, device := range devices {
		Device.Topics.Store(device.Code, device.Topic)
		Device.subscribeMqtt(&device)
	}
}

// 订阅mqtt
func (d *deviceService) subscribeMqtt(device *model.Device) {
	// 订阅mqtt的topic
	if device.Type == "cems" {
		global.MqttClient.Subscribe(device.Topic, 0, func(_ mqtt.Client, message mqtt.Message) {
			go ParseData.Cems(string(message.Payload()), message.Topic())
		})
	} else if device.Type == "dtu" {
		global.MqttClient.Subscribe(device.Topic, 0, func(_ mqtt.Client, message mqtt.Message) {
			go ParseData.Json(string(message.Payload()))
		})
	}
}

// 退订mqtt
func (d *deviceService) disSubscribeMqtt(device *model.Device) {
	global.MqttClient.Unsubscribe(device.Topic)
	d.Topics.Delete(device.Code)
}

func (d *deviceService) Add(device *model.Device) (any, error) {
	now := util.Now()
	device.CreateTime = now
	device.UpdateTime = now

	err := global.DeviceColl.FindOne(context.TODO(), bson.M{"code": device.Code}).Decode(&model.Device{})
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("code重复")
	}

	_, err = global.DeviceColl.InsertOne(context.Background(), device)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	sort.Ints(device.TimeInterval)
	if device.Type == "dtu" {
		for _, timeInterval := range device.TimeInterval {
			err = global.Database[timeInterval].CreateCollection(context.TODO(), device.Code)
			if err != nil {
				global.Log.Error(err.Error())
				return nil, err
			}
			_, err = global.Database[timeInterval].Collection(device.Code).Indexes().CreateOne(
				context.TODO(),
				mongo.IndexModel{
					Keys: bson.M{"createTime": -1},
				},
			)
			if err != nil {
				global.Log.Error(err.Error())
				return nil, err
			}
			global.CollMap[timeInterval].Store(device.Code, global.Database[timeInterval].Collection(device.Code))
		}
	}

	d.subscribeMqtt(device)
	return nil, nil
}

func (d *deviceService) PaginationFind(page *util.Page, name, Type, id string) (any, error) {
	opts := page.GetOpts()
	filter := bson.M{}
	if id != "" {
		oId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			global.Log.Error(err.Error())
			return nil, err
		}
		filter["group"] = oId
	}

	if name != "" {
		regexString := fmt.Sprintf(".*%v.*", name)
		filter["name"] = primitive.Regex{Pattern: regexString}
	}

	if Type != "" {
		filter["type"] = Type
	}

	cur, err := global.DeviceColl.Find(context.TODO(), filter, opts)
	var devices []model.Device
	err = cur.All(context.TODO(), &devices)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	count, err := global.DeviceColl.CountDocuments(context.TODO(), filter)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	cur, err = global.GroupColl.Find(context.TODO(), bson.M{})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	var groups []model.Group
	err = cur.All(context.TODO(), &groups)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	var resp []map[string]any
	for _, device := range devices {
		deviceMap := make(map[string]any)
		deviceMap["device"] = device
		if device.Group == primitive.NilObjectID {
			deviceMap["group"] = "默认分组"
		} else {
			for _, group := range groups {
				if device.Group == group.Id {
					deviceMap["group"] = group.Name
					break
				}
			}
		}

		resp = append(resp, deviceMap)
	}

	return gin.H{
		"devices": resp,
		"count":   count,
	}, nil
}

func (d *deviceService) Update(device *model.Device) (any, error) {
	var DBDevice model.Device
	device.UpdateTime = util.Now()
	filter := bson.M{"_id": device.Id}
	update := bson.M{"$set": device}
	err := global.DeviceColl.FindOne(context.TODO(), filter).Decode(&DBDevice)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	sort.Ints(device.TimeInterval)
	global.DeviceColl.FindOneAndUpdate(context.TODO(), filter, update)
	if !reflect.DeepEqual(device.TimeInterval, DBDevice.TimeInterval) {
		for _, timeInterval := range DBDevice.TimeInterval {
			global.CollMap[timeInterval].Delete(DBDevice.Code)
		}
		for _, timeInterval := range device.TimeInterval {
			err = global.Database[timeInterval].CreateCollection(context.TODO(), device.Code)
			if err != nil {
				global.Log.Error(err.Error())
				return nil, err
			}
			_, err = global.Database[timeInterval].Collection(device.Code).Indexes().CreateOne(
				context.TODO(),
				mongo.IndexModel{
					Keys: bson.M{"createTime": -1},
				},
			)
			global.CollMap[timeInterval].Store(device.Code, global.Database[timeInterval].Collection(device.Code))
		}
	}
	return nil, nil
}

func (d *deviceService) Delete(device *model.Device) (any, error) {
	filter := bson.M{
		"_id": device.Id,
	}
	_, err := global.DeviceColl.DeleteOne(context.TODO(), filter)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	if device.Type == "dtu" {
		opts := options.Count().SetLimit(1)
		for _, timeInterval := range device.TimeInterval {
			global.CollMap[timeInterval].Delete(device.Code)
			for _, database := range global.Database {
				collection := database.Collection(device.Code)
				count, err := collection.CountDocuments(context.TODO(), bson.M{}, opts)
				if err != nil {
					global.Log.Error(err.Error())
					return nil, err
				}
				if count == 0 {
					err = collection.Drop(context.TODO())
					if err != nil {
						global.Log.Error(err.Error())
						return nil, err
					}
				}
			}
		}
	} else if device.Type == "cems" {
		colls := []*mongo.Collection{
			global.CemsRealtimeColl,
			global.CemsMinuteColl,
			global.CemsHourColl,
			global.CemsDayColl,
		}

		filter := bson.M{
			"code": device.Code,
		}

		for _, coll := range colls {
			_, err := coll.DeleteMany(context.TODO(), filter)
			if err != nil {
				global.Log.Error(err.Error())
				return nil, err
			}
		}
	}
	d.disSubscribeMqtt(device)
	return nil, nil
}
