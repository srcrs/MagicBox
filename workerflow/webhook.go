package workerflow

import (
	"MagicBox/utils"
	"context"
	"net/url"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

// 通知
func (wf *WorkerFlowData) WebhookExecute(ctx context.Context, workflow, nodeId string) (interface{}, error) {
	// xxx := &LoopdataStrategy{}
	reqUrl := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.url`).String()
	method := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.method`).String()
	body := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.body`).String()
	contentType := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.contentType`).String()
	timeout := gjson.Get(workflow, `drawflow.nodes.#(id=="`+nodeId+`").data.timeout`).Int()

	headers := make(map[string]string)
	var reqParams interface{}

	body = utils.ReplaceAllVariable(body, wf.VariableMap)
	// 创建一个 JSON 对象
	bodyObj, err := simplejson.NewJson([]byte(body))
	if err != nil {
		utils.GLOBAL_LOGGER.Error("json Unmarshal: "+err.Error()+" body: "+body, zap.String("callid", ctx.Value("callid").(string)))
		return nil, err
	}
	if contentType == "form" {
		params := url.Values{}
		for k := range bodyObj.MustMap() {
			// default string
			if _, err := bodyObj.Get(k).String(); err == nil {
				params.Add(k, bodyObj.Get(k).MustString())
			} else {
				vv, err := bodyObj.Get(k).Encode()
				if err != nil {
					utils.GLOBAL_LOGGER.Error("body encode: " + err.Error() + " k:" + k)
				}
				params.Add(k, string(vv))
			}
		}
		reqParams = params
		utils.GLOBAL_LOGGER.Info("msg param", zap.Any("method: ", method), zap.Any("url: ", reqUrl), zap.Any("info: ", reqParams), zap.String("callid", ctx.Value("callid").(string)))
	}
	if response, err := utils.Request(method, reqUrl, headers, reqParams, time.Duration(timeout*int64(time.Second))); err != nil {
		utils.GLOBAL_LOGGER.Error("send msg: "+err.Error(), zap.String("callid", ctx.Value("callid").(string)))
	} else {
		utils.GLOBAL_LOGGER.Info("send msg rsp: "+response, zap.String("callid", ctx.Value("callid").(string)))
	}
	wf.NextNodeId = gjson.Get(workflow, `drawflow.edges.#(sourceHandle=="`+nodeId+`-output-fallback").target`).String()
	return nil, nil
}
