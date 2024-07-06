package workerflow

import (
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/chromedp/cdproto/cdp"
)

type WorkerFlowData struct {
	LoopDataElements map[string][]*cdp.Node
	VariableMap      *simplejson.Json
	NextNodeId       string
	ResetNextNodeId  string
	//为了兼容automa无法返回上一页问题
	CloseTabIsGoBack bool
}

func (wf *WorkerFlowData) SetVariableMap(variableName string, value interface{}) {
	variableName = "{{variables." + variableName + "}}"
	if strings.Contains(variableName, "$push:") {
		variableName = strings.ReplaceAll(variableName, "$push:", "")
		array := wf.VariableMap.Get(variableName).MustArray()
		array = append(array, value)
		wf.VariableMap.Set(variableName, array)
	} else {
		wf.VariableMap.Set(variableName, value)
	}
}
