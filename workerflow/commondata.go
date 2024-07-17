package workerflow

import (
	"MagicBox/utils"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/chromedp/cdproto/cdp"
	"go.uber.org/zap"
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
	variableNameNew := "{{variables." + variableName + "}}"
	variableNameOld := "{{variables@" + variableName + "}}"

	varList := []string{variableNameNew, variableNameOld}
	for _, variable := range varList {
		if strings.Contains(variable, "$push:") {
			variable = strings.ReplaceAll(variable, "$push:", "")
			array := wf.VariableMap.Get(variable).MustArray()
			array = append(array, value)
			wf.VariableMap.Set(variable, array)
		} else {
			wf.VariableMap.Set(variable, value)
		}
	}
}

func (wf *WorkerFlowData) GetVariableMapData(variable string) *simplejson.Json {
	utils.GLOBAL_LOGGER.Info("GetVariableMapData", zap.Any("variable", variable))
	if strings.Contains(variable, "$push:") {
		variable = strings.ReplaceAll(variable, "$push:", "")
	}
	if data, ok := wf.VariableMap.CheckGet(variable); ok {
		return data
	}
	return nil
}
