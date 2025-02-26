package global

import (
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

var (
	MongoClient      *mongo.Client
	MongoDatabase    *mongo.Database
	UserColl         *mongo.Collection
	RoleColl         *mongo.Collection
	ApiColl          *mongo.Collection
	LogColl          *mongo.Collection
	DeviceColl       *mongo.Collection
	LatestDataColl   *mongo.Collection
	CemsRealtimeColl *mongo.Collection
	CemsMinuteColl   *mongo.Collection
	CemsHourColl     *mongo.Collection
	CemsDayColl      *mongo.Collection
	MenuColl         *mongo.Collection
	GroupColl        *mongo.Collection
	CenterColl       *mongo.Collection
	ProjectColl      *mongo.Collection
	ReportColl       *mongo.Collection
	ParamColl        *mongo.Collection
	Database         [4]*mongo.Database //分别为1,15,60,60avg
	CollMap          [4]sync.Map
)
