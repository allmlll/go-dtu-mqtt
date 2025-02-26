package initiallize

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
)

func ApiCollInit(engine *gin.Engine) {
	routes := engine.Routes()
	_, err := global.ApiColl.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		global.Log.Error(err.Error())
	}
	var api model.Api
	for _, route := range routes {
		api.Url = route.Path
		api.Method = route.Method
		_, err = global.ApiColl.InsertOne(context.TODO(), api)
		if err != nil {
			global.Log.Error(err.Error())
		}
	}
}
