package workerflow

import (
	"MagicBox/utils"
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

//滑动页面到底部
func (wf *WorkerFlowData) DelayExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	sleepTime := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.time`).Int()
	if err := chromedp.Run(
		ctx,
		chromedp.Sleep(time.Duration(sleepTime)*time.Millisecond),
	); err != nil {
		utils.GLOBAL_LOGGER.Error("delay error: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
		return nil, err
	}

	return nil, nil
}
