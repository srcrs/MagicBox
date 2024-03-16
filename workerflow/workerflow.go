package workerflow

import (
	"MagicBox/utils"
	"context"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func AddCronTask(cronExpression, workflowId string) {
	utils.GLOBAL_LOGGER.Info(" 添加定时任务: " + workflowId)
	taskFunc := func() {
		Worker(workflowId)
	}
	_, err := utils.GLOBAL_CRON_WORKER.AddFunc(cronExpression, taskFunc)
	if err != nil {
		utils.GLOBAL_LOGGER.Error("cron add error", zap.Error(err))
	}
	taskFunc()
}

func Worker(workflowId string) {
	workflow := utils.GLOBAL_WORKFLOW_MAP[workflowId]
	nodesMap := utils.GetSortedEdges(workflow)

	//chromedp初始参数
	opts := utils.GetChromeConfigOpts()
	//创建一个上下文
	allCtx, cancel := chromedp.NewExecAllocator(
		context.Background(),
		opts...,
	)
	defer cancel()
	chromedpCtx, cancel := chromedp.NewContext(
		allCtx,
	)
	defer cancel()
	//创建超时时间
	chromedpCtx, cancel = context.WithTimeout(chromedpCtx, 100*time.Second)
	defer cancel()

	callId := uuid.New().String()
	chromedpCtx = context.WithValue(chromedpCtx, "callid", callId)

	workerflow := &WorkerFlowData{}

	workerflow.LoopDataElements = make(map[string][]*cdp.Node)
	workerflow.VariableMap = make(map[string]string)
	nextNodeId := gjson.Get(workflow, `drawflow.nodes.#(label=="trigger").id`).String()

	cookies := gjson.Get(workflow, `drawflow.nodes.#(label=="insert-data")#.data.dataList.#(name="cookies")`).String()
	cookies = gjson.Get(cookies, `0.value`).String()
	//判断是否存在cookie，加载到浏览器
	if cookies != "" {
		if err := chromedp.Run(
			chromedpCtx,
			utils.LoadCookies(`{"cookies":`+cookies+`}`),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("加载cookie出错", zap.Error(err), zap.String("callid", chromedpCtx.Value("callid").(string)))
		}
	}

	for len(nodesMap[nextNodeId]) > 0 {
		if len(nodesMap[nextNodeId]) == 1 {
			nextNodeId = nodesMap[nextNodeId][0]
		} else {
			nextNodeId = workerflow.NextNodeId
		}
		nodeId := nextNodeId
		nodeLabel := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").label`).String()
		switch nodeLabel {
		case "new-tab":
			workerflow.NewtabExecute(chromedpCtx, workflow, nodeId)
		case "loop-data":
			workerflow.LoopdataExecute(chromedpCtx, workflow, nodeId)
		case "get-text":
			workerflow.GettextExecute(chromedpCtx, workflow, nodeId)
		case "loop-breakpoint":
			delete(workerflow.LoopDataElements, nodeId)
		case "webhook":
			workerflow.WebhookExecute(chromedpCtx, workflow, nodeId)
		case "conditions":
			workerflow.ConditionsExecute(chromedpCtx, workflow, nodeId)
		case "event-click":
			workerflow.EventClickExecute(chromedpCtx, workflow, nodeId)
		case "insert-data":
			workerflow.InsertDataExecute(chromedpCtx, workflow, nodeId)
		default:
			utils.GLOBAL_LOGGER.Warn("break no label: "+nodeLabel, zap.String("callid", chromedpCtx.Value("callid").(string)))
		}
		if nodeLabel == "conditions" && nodeId == workerflow.NextNodeId {
			break
		}
	}
	//找出节点连接顺序
}
