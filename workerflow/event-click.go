package workerflow

import (
	"MagicBox/utils"
	"context"
	"regexp"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

// 点击页面元素
func (wf *WorkerFlowData) EventClickExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	//当前只有一个判断条件，先写死
	// waitSelectorTimeout := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.waitSelectorTimeout`).String()
	selector := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.selector`).String()
	utils.GLOBAL_LOGGER.Info("event click selector: "+selector, zap.String("callid", ctx.Value("callid").(string)))
	itemId := utils.GetVariableName(selector)
	if len(wf.LoopDataElements[itemId]) > 0 {
		re := regexp.MustCompile(`{{loopData[@.]{1}.*}}\s+`)
		selector = re.ReplaceAllString(selector, "")
		selector = utils.CssToXpath(selector)
		selector = wf.LoopDataElements[itemId][0].FullXPath() + selector
		utils.GLOBAL_LOGGER.Info("click element : "+selector, zap.String("callid", ctx.Value("callid").(string)))
		if err := chromedp.Run(
			ctx,
			chromedp.WaitVisible(selector),
			chromedp.Click(selector),
			chromedp.Sleep(2*time.Second),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("click element : "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
		}
	} else {
		if err := chromedp.Run(
			ctx,
			chromedp.Sleep(2*time.Second),
			chromedp.WaitVisible(selector),
			chromedp.Click(selector),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("event click error: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
			return nil, err
		}
	}

	return nil, nil
}
