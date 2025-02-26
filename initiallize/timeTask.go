package initiallize

import (
	"github.com/robfig/cron"
	"ruitong-new-service/global"
	"ruitong-new-service/service"

	"log"
)

func TimeTaskInit() {
	c := cron.New()

	// 添加一个 cron job，使用 cron 表达式
	//c.AddFunc("0 * * * * *", OneMin)
	//每15分钟
	err := c.AddFunc("0 */15 * * * *", FifteenMin)
	if err != nil {
		log.Fatal(err)
	}

	//每5分钟
	err = c.AddFunc("0 */5 * * * *", FiveMin)
	if err != nil {
		log.Fatal(err)
	}

	//每60分钟
	err = c.AddFunc("0 0 * * * *", OneHour)
	if err != nil {
		log.Fatal(err)
	}

	//每天
	err = c.AddFunc("0 3 0 * * *", OneDay)
	if err != nil {
		log.Fatal(err)
	}

	// 启动 cron scheduler
	c.Start()
	global.LogInfo.Info("cron start")

	// 阻塞主 goroutine，让程序持续运行
	select {}
}

func FiveMin() {
	service.LatestData.CheckStatus()
}

func FifteenMin() {
	service.TimeTask.FifMin()
	//service.TimeTask.FifMinCulCount()
}

func OneHour() {
	service.TimeTask.Hour()
	service.TimeTask.HourAvg()

}
func OneDay() {
	service.TimeTask.Day()
}
