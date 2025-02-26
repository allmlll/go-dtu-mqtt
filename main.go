package main

import (
	"ruitong-new-service/initiallize"
	"ruitong-new-service/router"
)

func init() {
	initiallize.Init()
}

func main() {
	//
	engine := router.GetEngine()
	//apiColl初始化
	initiallize.ApiCollInit(engine)
	if err := engine.Run(":8091"); err != nil {
		panic(err)
	}
}
