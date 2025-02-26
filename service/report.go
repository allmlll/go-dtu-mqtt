package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
	"ruitong-new-service/util"
)

type reportService struct {
}

var Report reportService

func (r *reportService) Get(page *util.Page, start, end, code string) (any, error) {

	filter := bson.M{
		"code": code,
	}
	if !(start == "" || end == "") {
		filter["createTime"] = bson.M{
			"$gte": start,
			"$lte": end,
		}
	}
	opt := page.GetOpts().SetSort(bson.M{"_id": -1})

	cur, err := global.ReportColl.Find(context.TODO(), filter, opt)
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
	count, err := global.ReportColl.CountDocuments(context.TODO(), filter)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	return gin.H{
		"data":  datas,
		"count": count,
	}, nil
}
