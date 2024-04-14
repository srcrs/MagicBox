package workerflow

import (
	"MagicBox/utils"
	"context"
	"fmt"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

//逻辑判断
func (wf *WorkerFlowData) ConditionsExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	//当前只有一个判断条件，先写死
	pathConditions := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.conditions`).Array()
	handle := ""
	for _, path := range pathConditions {
		id := gjson.Get(path.String(), `id`).String()
		childConditions := gjson.Get(path.String(), `conditions`).Array()
		isPathMatch := true
		for _, child := range childConditions {
			itemsConditions := gjson.Get(child.String(), `conditions`).Array()
			isChildMatch := true
			for _, itemData := range itemsConditions {
				items := gjson.Get(itemData.String(), `items`).String()
				compareResult := wf.compareValue(ctx, items)
				fmt.Println("result:", compareResult, items)
				if !compareResult {
					isChildMatch = false
					break
				}
			}
			if isChildMatch {
				isPathMatch = true
				break
			} else {
				isPathMatch = false
			}
		}
		if isPathMatch {
			handle = id
			break
		}
	}
	if handle == "" {
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
	utils.GLOBAL_LOGGER.Info("compareValue A: "+A+", B: "+B+", "+"compareType: "+typeItem, zap.String("callid", ctx.Value("callid").(string)))
	if typeItem == "cnt" {
		//包含
		if strings.Contains(A, B) {
			return true
		}
	} else if typeItem == "eq" {
		//相等
		if A == B {
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
		//先从variables中判断是否存在变量
		result = gjson.Get(data, "data.value").String()
		if wf.VariableMap[result] != "" {
			result = wf.VariableMap[result]
		}
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
	utils.GLOBAL_LOGGER.Info("nodeType: "+typeItem+", nodeValue: "+result, zap.String("callid", ctx.Value("callid").(string)))
	return result
}
