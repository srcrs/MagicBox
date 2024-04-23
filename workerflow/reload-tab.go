package workerflow

import (
	"MagicBox/utils"
	"context"

	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

//刷新当前页面
func (wf *WorkerFlowData) ReloadTabExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	utils.GLOBAL_LOGGER.Info("reload tab", zap.String("callid", ctx.Value("callid").(string)))

	if err := chromedp.Run(
		ctx,
		chromedp.Reload(),
	); err != nil {
		utils.GLOBAL_LOGGER.Error("reload tab: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
	}

	return nil, nil
}
