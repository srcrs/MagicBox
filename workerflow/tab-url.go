package workerflow

import (
	"MagicBox/utils"
	"context"

	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

//获取当前网页的url
func (wf *WorkerFlowData) TabUrlExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	// xxx := &LoopdataStrategy{}
	typeData := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.type`).String()
	variableName := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.variableName`).String()
	location := ""
	if typeData == "active-tab" {
		if err := chromedp.Run(
			ctx,
			chromedp.Location(&location),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("tab url error: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
			return nil, err
		}
		wf.VariableMap["{{variables@"+variableName+"}}"] = location
		utils.GLOBAL_LOGGER.Info("get url: "+variableName+", "+location, zap.String("callid", ctx.Value("callid").(string)))
	}
	return nil, nil
}
