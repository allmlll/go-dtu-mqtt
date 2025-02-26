package initiallize

func Init() {
	TimezoneInit()

	ViperInit()

	MongoInit()

	ZapInit()

	RedisInit()

	MqttInit()

	MinioInit()

	go TimeTaskInit()
}
