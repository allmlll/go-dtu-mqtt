package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
	"ruitong-new-service/util"
)

type logService struct {
}

var Log logService

func (l *logService) PaginationFind(page util.Page) (any, error) {
	opts := page.GetOpts().SetSort(bson.M{"_id": -1})

	cursor, err := global.LogColl.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	var logs []model.Log
	err = cursor.All(context.TODO(), &logs)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	count, err := global.LogColl.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	return gin.H{
		"logs":  logs,
		"count": count,
	}, nil
}

func (l *logService) DeleteAll() (any, error) {
	_, err := global.LogColl.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	return nil, nil
}
