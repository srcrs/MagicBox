package workerflow

import (
	"MagicBox/utils"
	"context"
	"regexp"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

// 页面表单填写
func (wf *WorkerFlowData) AttributeValueExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	selector := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.selector`).String()
	waitSelectorTimeout := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.waitSelectorTimeout`).Int()
	attributeName := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.attributeName`).String()
	variableName := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.variableName`).String()
	//是否开启等待
	waitForSelector := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.waitForSelector`).Bool()

	if waitForSelector {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(waitSelectorTimeout/1000*int64(time.Second)))
		defer cancel()
	}

	itemId := utils.GetVariableName(selector)
	if len(wf.LoopDataElements[itemId]) > 0 {
		re := regexp.MustCompile(`{{loopData[@.]{1}.*}}\s+`)
		selector = re.ReplaceAllString(selector, "")
		selector = utils.CssToXpath(selector)
		selector = wf.LoopDataElements[itemId][0].FullXPath() + selector
	}

	utils.GLOBAL_LOGGER.Info("AttributeValue selector: "+selector, zap.String("callid", ctx.Value("callid").(string)))

	result, flag := "", true

	if err := chromedp.Run(
		ctx,
		chromedp.WaitVisible(selector),
		chromedp.AttributeValue(selector, attributeName, &result, &flag),
	); err != nil {
		utils.GLOBAL_LOGGER.Error("AttributeValue: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
	}

	wf.SetVariableMap(variableName, result)

	return nil, nil
}
