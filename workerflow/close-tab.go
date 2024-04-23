package workerflow

import (
	"MagicBox/utils"
	"context"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

//刷新当前页面
func (wf *WorkerFlowData) CloseTabExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	utils.GLOBAL_LOGGER.Info("close tab", zap.String("callid", ctx.Value("callid").(string)))

	if err := chromedp.Run(
		ctx,
		page.Close(),
	); err != nil {
		utils.GLOBAL_LOGGER.Error("reload tab: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
	}

	return nil, nil
}
