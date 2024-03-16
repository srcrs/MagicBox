package workerflow

import "github.com/chromedp/cdproto/cdp"

type WorkerFlowData struct {
	LoopDataElements map[string][]*cdp.Node
	VariableMap      map[string]string
	NextNodeId       string
}
