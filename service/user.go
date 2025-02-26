package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"ruitong-new-service/global"
	"ruitong-new-service/model"
	"ruitong-new-service/util"
)

type userService struct {
}

var User userService

func (u *userService) Register(user *model.User) (any, error) {
	// 判断账号是否重复
	if err := global.UserColl.FindOne(context.TODO(), bson.M{"account": user.Account}).Decode(&model.User{}); !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("账号已存在")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	user.Password = string(hashedPassword)

	// 入库
	res, err := global.UserColl.InsertOne(context.TODO(), user)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	user.Id = res.InsertedID.(primitive.ObjectID)

	// 生成token
	tokenString, err := util.CreateToken(user.Id)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	// 查找角色
	count, err := global.RoleColl.CountDocuments(context.TODO(), bson.M{"_id": user.Role})
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("角色不存在")
	}

	return gin.H{
		"user":  user,
		"token": util.TokenPrefix + " " + tokenString,
	}, nil
}

func (u *userService) Login(account, password string) (any, error) {
	var DBUser model.User

	// 根据账号找用户
	err := global.UserColl.FindOne(context.TODO(), bson.M{"account": account}).Decode(&DBUser)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("账号未注册")
	} else if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	// 判断密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(DBUser.Password), []byte(password))
	if err != nil {
		global.Log.Error(err.Error())
		return nil, errors.New("密码错误")
	}

	// 生成token
	tokenString, err := util.CreateToken(DBUser.Id)
	if err != nil {
		return nil, err
	}

	// 查找角色
	var role model.Role
	err = global.RoleColl.FindOne(context.TODO(), bson.M{"_id": DBUser.Role}).Decode(&role)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	return gin.H{
		"user":  DBUser,
		"token": util.TokenPrefix + " " + tokenString,
		"role":  role,
	}, nil
}

func (u *userService) Get(page util.Page, name string) (any, error) {
	opts := page.GetOpts()
	var filter bson.M
	if name != "" {
		regexString := fmt.Sprintf(".*%v.*", name)
		filter = bson.M{
			"name": primitive.Regex{Pattern: regexString},
		}
	}

	cursor, err := global.UserColl.Find(context.TODO(), filter, opts)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	var users []model.User
	err = cursor.All(context.TODO(), &users)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	count, err := global.UserColl.CountDocuments(context.TODO(), filter)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	return gin.H{
		"users": users,
		"count": count,
	}, nil
}

func (u *userService) Update(user *model.User) (any, error) {
	filter := bson.M{"_id": user.Id}
	update := bson.M{"$set": bson.M{"avatar": user.Avatar, "account": user.Account,
		"phone": user.Phone, "name": user.Name, "sex": user.Sex, "role": user.Role}}
	global.UserColl.FindOneAndUpdate(context.TODO(), filter, update)

	return user, nil
}

func (u *userService) Delete(id primitive.ObjectID) (any, error) {
	filter := bson.M{
		"_id": id,
	}
	_, err := global.UserColl.DeleteOne(context.TODO(), filter)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	return nil, nil
}

func (u *userService) ResetPass(id primitive.ObjectID, phone string) (any, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(phone), bcrypt.DefaultCost)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	password := string(hashedPassword)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"password": password}}
	global.UserColl.FindOneAndUpdate(context.TODO(), filter, update)

	return nil, nil
}

func (u *userService) ChangePass(id any, newPass, oldPass string) (any, error) {
	var DBUser model.User
	filter := bson.M{"_id": id}

	// 根据账号找用户
	err := global.UserColl.FindOne(context.TODO(), filter).Decode(&DBUser)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	// 判断密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(DBUser.Password), []byte(oldPass))
	if err != nil {
		global.Log.Error(err.Error())
		return nil, errors.New("原密码错误")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	password := string(hashedPassword)
	update := bson.M{"$set": bson.M{"password": password}}
	global.UserColl.FindOneAndUpdate(context.TODO(), filter, update)

	return nil, nil
}
