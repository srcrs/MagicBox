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
func (wf *WorkerFlowData) LoopdataExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	//loop data类型
	loopThrough := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.loopThrough`).String()
	//css路径
	elementSelector := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.elementSelector`).String()
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
	if loopThrough == "elements" {
		if err := chromedp.Run(
			ctx,
			chromedp.WaitVisible(elementSelector),
			chromedp.Nodes(elementSelector, &nodeElements),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("loop data error: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
			return nil, err
		}
		wf.LoopDataElements[loopId] = nodeElements
	}
	return nil, nil
}
