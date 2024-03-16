package workerflow

import (
	"MagicBox/utils"
	"context"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

//逻辑判断
func (wf *WorkerFlowData) ConditionsExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	//当前只有一个判断条件，先写死
	id := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.conditions.0.id`).String()
	items := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.conditions.0.conditions.0.conditions.0.items`).String()
	compareResult := wf.compareValue(ctx, items)
	handle := ""
	if compareResult {
		handle = id
	} else {
		handle = "fallback"
	}
	wf.NextNodeId = gjson.Get(workflow, `drawflow.edges.#(sourceHandle=="`+nodeId+"-output-"+handle+`").target`).String()
	return nil, nil
}

//判断节点操作类型
func (wf *WorkerFlowData) compareValue(ctx context.Context, items string) bool {
	A := wf.getValue(ctx, gjson.Get(items, "0").String())
	B := wf.getValue(ctx, gjson.Get(items, "2").String())
	typeItem := gjson.Get(items, "1.type").String()
	if typeItem == "cnt" {
		//包含
		if strings.Contains(A, B) {
			return true
		}
	}
	return false
}

//获取节点的真实value
func (wf *WorkerFlowData) getValue(ctx context.Context, data string) string {
	result := ""
	typeItem := gjson.Get(data, "type").String()

	if typeItem == "element#text" {
		selector := gjson.Get(data, "data.selector").String()
		//从当前网页获取元素的value
		if err := chromedp.Run(
			ctx,
			chromedp.WaitVisible(selector),
			chromedp.TextContent(selector, &result),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("Conditions getValue error: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
		}
	} else if typeItem == "value" {
		result = gjson.Get(data, "data.value").String()
	} else if typeItem == "element#attribute" {
		selector := gjson.Get(data, "data.selector").String()
		attrName := gjson.Get(data, "data.attrName").String()
		flag := true
		if err := chromedp.Run(
			ctx,
			chromedp.WaitVisible(selector),
			chromedp.AttributeValue(selector, attrName, &result, &flag),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("Conditions getValue error: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
		}

	}
	return result
}
