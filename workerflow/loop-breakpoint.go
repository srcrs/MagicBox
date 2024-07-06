package workerflow

import (
	"MagicBox/utils"
	"context"

	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

// 循环一次结束
func (wf *WorkerFlowData) LoopBreakPointExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	loopId := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.loopId`).String()
	if len(wf.LoopDataElements[loopId]) > 1 {
		//移除第一个元素
		wf.LoopDataElements[loopId] = wf.LoopDataElements[loopId][1:]
		//将节点重新置到循环开始的位置
		loopElements := gjson.Get(workflow, `drawflow.nodes.#(label==loop-elements)`).Array()
		for _, element := range loopElements {
			loopIdTmp := gjson.Get(element.String(), `data.loopId`).String()
			if loopIdTmp == loopId {
				wf.ResetNextNodeId = gjson.Get(element.String(), `id`).String()
				utils.GLOBAL_LOGGER.Info("reset Node: "+wf.ResetNextNodeId, zap.Any("xxx", len(wf.LoopDataElements[loopId])))
				break
			}
		}
	}
	return nil, nil
}
