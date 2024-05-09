package workerflow

import (
	"github.com/chromedp/cdproto/cdp"
)

type WorkerFlowData struct {
	LoopDataElements map[string][]*cdp.Node
	VariableMap      map[string]string
	NextNodeId       string
	ResetNextNodeId  string
	//为了兼容automa无法返回上一页问题
	CloseTabIsGoBack bool
}
