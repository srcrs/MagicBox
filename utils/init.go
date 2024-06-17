package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

//初始化定时任务
func InitCronWorker() {
	GLOBAL_CRON_WORKER = cron.New(
		// cron.WithSeconds(),
		cron.WithChain(cron.Recover(cron.DefaultLogger)),
		//当前如果有任务在执行则加入到任务队列中
		cron.WithChain(cron.DelayIfStillRunning(cron.DefaultLogger)),
	)
	GLOBAL_CRON_WORKER.Start()
}

//初始化日志
func InitLogger() {
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleOutput := zapcore.Lock(zapcore.AddSync(os.Stdout))
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileOutput := zapcore.Lock(zapcore.AddSync(&lumberjack.Logger{
		Filename:   "MagicBox.log",
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     30, //days
	}))
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleOutput, zap.InfoLevel),
		zapcore.NewCore(fileEncoder, fileOutput, zap.InfoLevel),
	)
	GLOBAL_LOGGER = zap.New(core)
	GLOBAL_LOGGER = GLOBAL_LOGGER.WithOptions(zap.AddCaller())
}

//初始化配置
func InitConfig(root string) {
	GLOBAL_WORKFLOW_MAP = make(map[string]string)
	GLOBAL_WROK_CRONING = make(map[string]cron.EntryID)
	var files []os.FileInfo

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, info)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		filePath := root + file.Name()
		// 读取文件内容
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		GLOBAL_WORKFLOW_MAP[GenerateMD5(filePath)] = string(fileContent)
	}
}
