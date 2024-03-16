package workerflow

import (
	"MagicBox/utils"
	"context"
	"fmt"
	"regexp"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func (wf *WorkerFlowData) GettextExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	// xxx := &LoopdataStrategy{}
	// findBy := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.findBy`).String()
	// waitForSelector := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.waitForSelector`).Bool()
	// waitSelectorTimeout := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.waitSelectorTimeout`).Int()
	selector := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.selector`).String()
	variableName := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.variableName`).String()
	// dataColumn := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.dataColumn`).String()
	itemId := utils.GetVariableName(selector)
	if len(wf.LoopDataElements[itemId]) > 0 {
		re := regexp.MustCompile(`{{loopData@.*}}\s+`)
		selector = re.ReplaceAllString(selector, "")

		var chrnodes []*cdp.Node

		for _, node := range wf.LoopDataElements[itemId] {
			chromedp.Run(
				ctx,
				chromedp.Nodes(selector, &chrnodes, chromedp.FromNode(node)),
			)
		}
		for _, node := range chrnodes {
			all := ""
			chromedp.Run(
				ctx,
				chromedp.TextContent(node.FullXPath(), &all),
			)
			fmt.Println(all)
		}
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
