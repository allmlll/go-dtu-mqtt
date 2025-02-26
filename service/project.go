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

type projectService struct {
}

var Project projectService

func (p *projectService) Update(project *model.Project) (any, error) {
	_, err := global.ProjectColl.UpdateOne(context.TODO(), bson.M{"_id": project.Id}, bson.M{"$set": project})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	return nil, nil
}

func (p *projectService) Get(page *util.Page, url string) (any, error) {
	opts := page.GetOpts()
	var projects []model.Project
	var filter bson.M
	if url != "" {
		regexString := fmt.Sprintf(".*%v.*", url)
		filter = bson.M{"url": primitive.Regex{Pattern: regexString}}
	}

	cur, err := global.ProjectColl.Find(context.TODO(), filter, opts)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	err = cur.All(context.TODO(), &projects)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	count, err := global.ProjectColl.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"projects": projects,
		"count":    count,
	}, nil
}

func (p *projectService) Add(project *model.Project) (any, error) {
	_, err := global.ProjectColl.InsertOne(context.TODO(), project)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	return nil, nil
}

func (p *projectService) Delete(project *model.Project) (any, error) {
	_, err := global.ProjectColl.DeleteOne(context.TODO(), bson.M{"_id": project.Id})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	return nil, nil
}
