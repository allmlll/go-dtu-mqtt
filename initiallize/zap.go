package initiallize

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"ruitong-new-service/global"
	"ruitong-new-service/model"

	"os"
	"runtime"
	"time"
)

type mongoWriter struct {
}

var mongodb mongoWriter

func ZapInit() {
	// 创建 encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	console := zapcore.Lock(os.Stdout)
	var core zapcore.Core
	if runtime.GOOS == "linux" {
		core = zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(&mongodb), zap.NewAtomicLevel())
	} else {
		core = zapcore.NewCore(encoder, console, zap.NewAtomicLevel())
	}
	global.LogInfo = zap.New(zapcore.NewCore(encoder, console, zap.NewAtomicLevel()), zap.AddCaller())
	// 创建 logger
	global.Log = zap.New(core, zap.AddCaller())
	defer global.Log.Sync()
}

func (mw *mongoWriter) Write(p []byte) (n int, err error) {
	var logMap map[string]interface{}
	if err = json.Unmarshal(p, &logMap); err != nil {
		return 0, err
	}
	log := model.Log{
		Level:  logMap["level"].(string),
		Time:   time.Now().Format("2006-01-02 15:04:05"),
		Caller: logMap["caller"].(string),
		Msg:    logMap["msg"].(string),
	}
	_, err = global.LogColl.InsertOne(context.TODO(), log)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}
