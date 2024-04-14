package workerflow

import (
	"MagicBox/utils"
	"context"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

//遍历页面列表数据
func (wf *WorkerFlowData) LoopElementsExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	//css/xpath路径
	selector := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.selector`).String()
	//是否开启等待
	waitForSelector := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.waitForSelector`).Bool()

	waitSelectorTimeout := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.waitSelectorTimeout`).Int()

	//存储loopId
	loopId := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.loopId`).String()

	if waitForSelector {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(waitSelectorTimeout/1000*int64(time.Second)))
		defer cancel()
	}
	var nodeElements []*cdp.Node
	//代表获取elements
	if err := chromedp.Run(
		ctx,
		chromedp.WaitVisible(selector),
		chromedp.Nodes(selector, &nodeElements),
	); err != nil {
		utils.GLOBAL_LOGGER.Error("new tab error: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
		return nil, err
	}
	wf.LoopDataElements[loopId] = nodeElements
	return nil, nil
}
