package workerflow

import (
	"context"

	"github.com/tidwall/gjson"
)

func (wf *WorkerFlowData) InsertDataExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	dataList := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.dataList`).Array()
	for _, result := range dataList {
		typeData := gjson.Get(result.String(), "type").String()
		name := gjson.Get(result.String(), "name").String()
		value := gjson.Get(result.String(), "value").String()
		if typeData == "variable" {
			wf.SetVariableMap(name, value)
		}
	}
	return nil, nil
}
