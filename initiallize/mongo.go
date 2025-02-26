package initiallize

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
)

func MongoInit() {
	account := viper.GetString("mongo.account")
	password := viper.GetString("mongo.password")
	address := viper.GetString("mongo.address")
	port := viper.GetString("mongo.port")
	database := viper.GetString("mongo.database")
	url := "mongodb://" + account + ":" + password + "@" + address + ":" + port
	//mongo初始化
	if global.MongoClient == nil {
		global.MongoClient = getMongoClient(url)
	}

	//hbrt数据库及coll初始化
	global.MongoDatabase = global.MongoClient.Database(database)
	{
		global.DeviceColl = global.MongoDatabase.Collection("device")
		global.UserColl = global.MongoDatabase.Collection("user")
		global.LogColl = global.MongoDatabase.Collection("log")
		global.ApiColl = global.MongoDatabase.Collection("api")
		global.RoleColl = global.MongoDatabase.Collection("role")
		global.LatestDataColl = global.MongoDatabase.Collection("latestData")
		global.CemsRealtimeColl = global.MongoDatabase.Collection("cems-realtime")
		global.CemsMinuteColl = global.MongoDatabase.Collection("cems-minute")
		global.CemsHourColl = global.MongoDatabase.Collection("cems-hour")
		global.CemsDayColl = global.MongoDatabase.Collection("cems-day")
		global.MenuColl = global.MongoDatabase.Collection("menu")
		global.GroupColl = global.MongoDatabase.Collection("group")
		global.CenterColl = global.MongoDatabase.Collection("center")
		global.ProjectColl = global.MongoDatabase.Collection("project")
		global.ReportColl = global.MongoDatabase.Collection("report")
		global.ParamColl = global.MongoDatabase.Collection("param")
	}
	global.Database[0] = global.MongoClient.Database("1")
	global.Database[1] = global.MongoClient.Database("15")
	global.Database[2] = global.MongoClient.Database("60")
	global.Database[3] = global.MongoClient.Database("60avg")

	var devices []model.Device
	cur, err := global.DeviceColl.Find(context.TODO(), bson.M{})
	err = cur.All(context.TODO(), &devices)
	if err != nil {
		fmt.Println(err)
	}
	for _, device := range devices {
		if device.Type == "dtu" {
			for _, timeInterval := range device.TimeInterval {
				global.CollMap[timeInterval].Store(device.Code, global.Database[timeInterval].Collection(device.Code))
			}
			//for i := range global.CollMap {
			//	global.CollMap[i].Store(device.Code, global.Database[i].Collection(device.Code))
			//}
		}
	}
}

func getMongoClient(uri string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(uri)
	// 使用客户端选项和上下文（context.TODO()）连接到 MongoDB 数据库
	MongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println(err)
	}
	// 使用 Ping 方法检测与数据库的连接状态
	if err = MongoClient.Ping(context.TODO(), nil); err != nil {
		log.Println(err)
	}
	fmt.Println("mongodb连接成功")
	return MongoClient
}
