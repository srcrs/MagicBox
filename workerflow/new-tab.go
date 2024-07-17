package workerflow

import (
	"MagicBox/utils"
	"context"

	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

// 打开网页
func (wf *WorkerFlowData) NewtabExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	// xxx := &LoopdataStrategy{}
	reqUrl := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.url`).String()
	if data := wf.GetVariableMapData(reqUrl); data != nil {
		reqUrl = data.MustString()
	}
	if reqUrlTmp, err := utils.FixURL(reqUrl); err == nil {
		reqUrl = reqUrlTmp
	} else {
		utils.GLOBAL_LOGGER.Error("new tab fixUrl: "+err.Error(), zap.Any("url", reqUrl), zap.String("callid", ctx.Value("callid").(string)))
		return nil, err
	}
	if err := chromedp.Run(
		ctx,
		chromedp.Navigate(reqUrl),
	); err != nil {
		utils.GLOBAL_LOGGER.Error("new tab error: "+err.Error(), zap.Any("url", reqUrl), zap.String("callid", ctx.Value("callid").(string)))
		return nil, err
	}
	return nil, nil
}
