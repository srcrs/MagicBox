package workerflow

import (
	"MagicBox/utils"
	"context"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

// 关闭当前页面
func (wf *WorkerFlowData) CloseTabExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	utils.GLOBAL_LOGGER.Info("close tab", zap.String("callid", ctx.Value("callid").(string)))
	if wf.CloseTabIsGoBack {
		if err := chromedp.Run(
			ctx,
			chromedp.Sleep(3*time.Second),
			chromedp.Evaluate(`history.go(-1)`, nil),
			chromedp.Sleep(3*time.Second),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("close tab revise goBackPage: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
		}
		wf.CloseTabIsGoBack = false
	} else {
		if err := chromedp.Run(
			ctx,
			page.Close(),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("close tab: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
		}
	}

	return nil, nil
}
