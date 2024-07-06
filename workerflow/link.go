package workerflow

import (
	"MagicBox/utils"
	"context"
	"regexp"

	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

// open link
func (wf *WorkerFlowData) LinkExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	selector := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.selector`).String()
	// openInNewTab := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.openInNewTab`).String()
	itemId := utils.GetVariableName(selector)
	if len(wf.LoopDataElements[itemId]) > 0 {
		re := regexp.MustCompile(`{{loopData[@.]{1}.*}}\s+`)
		selector = re.ReplaceAllString(selector, "")
		selector = utils.CssToXpath(selector)
		selector = wf.LoopDataElements[itemId][0].FullXPath() + selector
	}
	utils.GLOBAL_LOGGER.Info("Link selector: "+selector, zap.String("callid", ctx.Value("callid").(string)))

	if err := chromedp.Run(
		ctx,
		chromedp.EvaluateAsDevTools(`var link = document.evaluate("`+selector+`", document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue; link.setAttribute('target', '_self');`, nil),
		chromedp.Click(selector),
	); err != nil {
		utils.GLOBAL_LOGGER.Error("open link url error: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
	}
	wf.CloseTabIsGoBack = true
	return nil, nil
}
