package workerflow

import (
	"MagicBox/utils"
	"context"

	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

//点击页面元素
func (wf *WorkerFlowData) EventClickExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	//当前只有一个判断条件，先写死
	// waitSelectorTimeout := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.waitSelectorTimeout`).String()
	selector := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.selector`).String()
	if err := chromedp.Run(
		ctx,
		chromedp.WaitVisible(selector),
		chromedp.Click(selector),
	); err != nil {
		utils.GLOBAL_LOGGER.Error("event click error: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
		return nil, err
	}

	return nil, nil
}
