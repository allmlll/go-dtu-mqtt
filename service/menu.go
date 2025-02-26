package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
)

type menuService struct{}

var Menu menuService

func (m *menuService) GetMenu(id any) (any, error) {
	var user model.User
	err := global.UserColl.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	var role model.Role
	err = global.RoleColl.FindOne(context.TODO(), bson.M{"_id": user.Role}).Decode(&role)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	var centerCascades []model.Cascade
	opts := options.Find().SetSort(bson.M{"index": 1})
	cur, err := global.CenterColl.Find(context.TODO(), bson.M{"_id": bson.M{"$in": role.Centers}}, opts)
	err = cur.All(context.TODO(), &centerCascades)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	for i, centerCascade := range centerCascades {
		cur, err = global.GroupColl.Find(context.TODO(), bson.M{"center": centerCascade.Value}, opts)
		var groupCascades []model.Cascade
		var newGroupCascades []model.Cascade
		err = cur.All(context.TODO(), &groupCascades)
		if err != nil {
			global.Log.Error(err.Error())
			return nil, err
		}
		for _, groupCascade := range groupCascades {
			cur, err = global.DeviceColl.Find(context.TODO(), bson.M{"group": groupCascade.Value, "code": bson.M{"$in": role.Codes}}, opts)
			var deviceCascades []model.Cascade
			err = cur.All(context.TODO(), &deviceCascades)
			if err != nil {
				global.Log.Error(err.Error())
				return nil, err
			}
			if len(deviceCascades) != 0 {
				groupCascade.Children = deviceCascades
				newGroupCascades = append(newGroupCascades, groupCascade)
			}
		}
		centerCascades[i].Children = newGroupCascades
	}

	return centerCascades, nil
}
