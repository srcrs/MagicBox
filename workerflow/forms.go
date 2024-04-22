package workerflow

import (
	"MagicBox/utils"
	"context"

	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

//页面表单填写
func (wf *WorkerFlowData) FormsExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	selector := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.selector`).String()
	itemType := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.type`).String()
	value := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.value`).String()
	utils.GLOBAL_LOGGER.Info("forms selector: "+selector, zap.String("callid", ctx.Value("callid").(string)))

	if itemType == "text-field" {
		utils.GLOBAL_LOGGER.Info("forms set value: "+value, zap.String("callid", ctx.Value("callid").(string)))
		if err := chromedp.Run(
			ctx,
			chromedp.WaitVisible(selector),
			chromedp.SetValue(selector, value),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("forms set value: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
		}
	} else {
		utils.GLOBAL_LOGGER.Error("forms need coder: ", zap.String("callid", ctx.Value("callid").(string)))
	}

	return nil, nil
}
