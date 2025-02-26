package initiallize

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"log"
	"ruitong-new-service/global"
	"ruitong-new-service/service"

	"time"
)

func MqttInit() {
	id := viper.GetString("mqtt.id")
	mqttClientInit(id, &global.MqttClient)
	service.Device.TopicInit()
	go service.Param.ParamPublish()
}

func mqttClientInit(id string, client *mqtt.Client) {
	//配置
	address := viper.GetString("mqtt.address")
	port := viper.GetString("mqtt.port")
	url := "tcp://" + address + ":" + port
	account := viper.GetString("mqtt.account")
	password := viper.GetString("mqtt.password")
	//设置地址用户名和密码
	clientOptions := mqtt.NewClientOptions().AddBroker(url).SetUsername(account).SetPassword(password)
	//设置客户端id
	clientOptions.SetClientID(id)
	//掉线后不清除session
	clientOptions.SetCleanSession(false)
	//设置自动重连
	clientOptions.SetAutoReconnect(true)
	//设置处理函数
	clientOptions.OnConnect = func(client mqtt.Client) {
		global.LogInfo.Info("mqttClient," + id + "与服务端建立连接成功！")
	}
	clientOptions.OnConnectionLost = func(client mqtt.Client, err error) {
		global.Log.Error("mqttClient," + id + "与服务端断开连接: " + err.Error())
	}
	clientOptions.OnReconnecting = func(client mqtt.Client, options *mqtt.ClientOptions) {
		log.Println("mqttClient," + id + "与服务端重连成功！")
	}

	*client = mqtt.NewClient(clientOptions) //客户端建立
	//客户端连接判断
	if token := (*client).Connect(); token.WaitTimeout(time.Duration(60)*time.Second) && token.Wait() && token.Error() != nil {
		global.Log.Error(token.Error().Error() + "mqtt")
	}
}
