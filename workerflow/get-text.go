package workerflow

import (
	"MagicBox/utils"
	"context"
	"regexp"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func (wf *WorkerFlowData) GettextExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	selector := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.selector`).String()
	variableName := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.variableName`).String()
	// findBy := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.findBy`).String()
	itemId := utils.GetVariableName(selector)
	if len(wf.LoopDataElements[itemId]) > 0 {
		re := regexp.MustCompile(`{{loopData@.*}}\s+`)
		selector = re.ReplaceAllString(selector, "")
		selector = utils.CssToXpath(selector)
		selector = wf.LoopDataElements[itemId][0].FullXPath() + selector
		utils.GLOBAL_LOGGER.Info("get text selector: "+selector, zap.String("callid", ctx.Value("callid").(string)))
		all := ""
		if err := chromedp.Run(
			ctx,

			chromedp.WaitVisible(selector),
			chromedp.TextContent(selector, &all),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("get text : "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
		}
		all = strings.Replace(all, "\n", "", -1)
		all = strings.Replace(all, "\t", "", -1)
		utils.GLOBAL_LOGGER.Info("get text : "+all+" selector: "+wf.LoopDataElements[itemId][0].FullXPath()+" "+selector, zap.String("callid", ctx.Value("callid").(string)))
		wf.VariableMap["{{variables@"+variableName+"}}"] = all
	} else {
		textContent := ""
		if err := chromedp.Run(
			ctx,
			chromedp.TextContent(selector, &textContent),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("get text error: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
			return nil, err
		}
		wf.VariableMap["{{variables@"+variableName+"}}"] = textContent
	}
	return nil, nil
}
