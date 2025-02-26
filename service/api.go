package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
	"ruitong-new-service/util"
)

type apiService struct {
}

var Api apiService

func (a *apiService) PaginationFind(page util.Page, url, method string) (any, error) {
	opts := page.GetOpts()
	var filter bson.M
	if method != "" {
		if url != "" {
			regexString := fmt.Sprintf(".*%v.*", url)
			filter = bson.M{
				"url":    primitive.Regex{Pattern: regexString},
				"method": method,
			}
		} else {
			filter = bson.M{
				"method": method,
			}
		}
	} else {
		if url != "" {
			regexString := fmt.Sprintf(".*%v.*", url)
			filter = bson.M{
				"url": primitive.Regex{Pattern: regexString},
			}
		}
	}

	cursor, err := global.ApiColl.Find(context.TODO(), filter, opts)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	var apis []model.Api
	err = cursor.All(context.TODO(), &apis)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	count, err := global.ApiColl.CountDocuments(context.TODO(), filter)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	return gin.H{
		"apis":  apis,
		"count": count,
	}, nil
}

func (a *apiService) Get() (any, error) {

	cursor, err := global.ApiColl.Find(context.TODO(), bson.M{})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	var apis []model.Api
	err = cursor.All(context.TODO(), &apis)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	count, err := global.ApiColl.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	return gin.H{
		"apis":  apis,
		"count": count,
	}, nil
}
