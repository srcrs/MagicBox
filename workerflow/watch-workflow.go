package workerflow

import (
	"MagicBox/utils"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

// MonitorFolder 监控文件夹中的文件变动
func WatchConfigChanges(folderPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.GLOBAL_LOGGER.Fatal("WatchConfigChanges: " + err.Error())
		log.Fatal(err)
	}
	defer watcher.Close()

	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			utils.GLOBAL_LOGGER.Error("Error walking config folder: " + err.Error())
			return err
		}
		if info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				utils.GLOBAL_LOGGER.Error("Error adding watch for file: " + err.Error())
			}
		}
		return nil
	})

	if err != nil {
		utils.GLOBAL_LOGGER.Fatal("Error: " + err.Error())
	}

	utils.GLOBAL_LOGGER.Info("Start WatchConfigChanges: " + folderPath)

	for {
		time.Sleep(time.Second)
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			utils.GLOBAL_LOGGER.Info("file update: " + event.Name)
			if isConfigFile(event.Name) {
				// 当配置文件发生更改时，重新加载配置
				reloadConfigFunc(event.Name)
			}
		case err, ok := <-watcher.Errors:
			utils.GLOBAL_LOGGER.Error("err: " + err.Error())
			if !ok {
				return
			}
		default:
			continue
		}

	}
}

// isConfigFile 检查给定的文件路径是否是配置文件
func isConfigFile(filePath string) bool {
	// 这里可以根据需要自定义配置文件的判断逻辑，
	// 例如根据文件名后缀或其他文件内容特征进行判断
	extension := strings.ToLower(filepath.Ext(filePath))
	return extension == ".json"
}

//配置文件热更新
func reloadConfigFunc(filePath string) {
	workflowId := utils.GenerateMD5(filePath)
	utils.GLOBAL_LOGGER.Info("Reloading configuration: " + filePath)
	if cronId, ok := utils.GLOBAL_WROK_CRONING[workflowId]; ok {
		utils.GLOBAL_CRON_WORKER.Remove(cronId)
		//移除任务Id
		utils.GLOBAL_LOGGER.Info("remove cron task: ", zap.Any("cronID", cronId))
	}
	// 读取文件内容
	if fileContent, err := ioutil.ReadFile(filePath); err == nil {
		utils.GLOBAL_WORKFLOW_MAP[workflowId] = string(fileContent)
		cronExpression := gjson.Get(string(fileContent), `drawflow.nodes.#(label=="trigger").data.triggers.#(type="cron-job").data.expression`)
		utils.GLOBAL_LOGGER.Info("update cron: " + workflowId)
		taskFunc := func() {
			Worker(workflowId)
		}
		cronId, err := utils.GLOBAL_CRON_WORKER.AddFunc(cronExpression.String(), taskFunc)
		if err != nil {
			utils.GLOBAL_LOGGER.Error("cron add error", zap.Error(err))
		}
		utils.GLOBAL_WROK_CRONING[workflowId] = cronId
		taskFunc()
	}
}
