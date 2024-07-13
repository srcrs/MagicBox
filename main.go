package main

import (
	"MagicBox/cmd"
	"MagicBox/utils"
	"MagicBox/workerflow"
	"os"

	"github.com/tidwall/gjson"
)

func main() {
	//初始化日志
	utils.InitLogger()
	//初始化配置
	configRoot := "configs/"
	utils.InitConfig(configRoot)

	//接收外部传输参数
	if len(os.Args) > 1 {
		cmd.Execute()
		return
	}

	//配置文件热更新
	go workerflow.WatchConfigChanges(configRoot)
	//初始化定时任务
	utils.InitCronWorker()
	utils.GLOBAL_LOGGER.Info("当前版本: 2.2.2")
	//任务执行
	for k, v := range utils.GLOBAL_WORKFLOW_MAP {
		//使用gjson解析
		cronExpression := gjson.Get(v, `drawflow.nodes.#(label=="trigger").data.triggers.#(type="cron-job").data.expression`)
		workerflow.AddCronTask(cronExpression.String(), k)
	}
	select {}
}
