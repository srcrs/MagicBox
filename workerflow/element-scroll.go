package workerflow

import (
	"MagicBox/utils"
	"context"

	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

//滑动页面到底部
func (wf *WorkerFlowData) ElementScrollExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	// selector := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.selector`).String()
	// waitSelectorTimeout := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.waitSelectorTimeout`).Int()
	// chromdpCtx, cancel := context.WithTimeout(ctx, time.Duration(waitSelectorTimeout)*time.Millisecond)
	// defer cancel()
	if err := chromedp.Run(
		ctx,
		// chromedp.WaitVisible(selector),
		chromedp.EvaluateAsDevTools(`window.scrollTo(0, document.body.scrollHeight);`, nil),
	); err != nil {
		utils.GLOBAL_LOGGER.Error("element-scroll error: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
		return nil, err
	}

	return nil, nil
}
