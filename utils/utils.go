package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/tidwall/gjson"
)

func GenerateMD5(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

func Request(method, requrl string, headers map[string]string, data interface{}, timeout time.Duration) (string, error) {
	client := &http.Client{Timeout: timeout}

	var req *http.Request
	var err error

	switch d := data.(type) {
	case nil:
		req, err = http.NewRequest(method, requrl, nil)
		if err != nil {
			return "", err
		}
	case url.Values:
		req, err = http.NewRequest(method, requrl, strings.NewReader(d.Encode()))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	default:
		jsonData, err := json.Marshal(d)
		if err != nil {
			return "", err
		}
		req, err = http.NewRequest(method, requrl, bytes.NewBuffer(jsonData))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf := new(strings.Builder)

	return buf.String(), nil
}

// 获取节点连接后的顺序
func GetSortedEdges(workflow string) map[string][]string {
	edges := gjson.Get(workflow, `drawflow.edges`).Array()
	sourceToTargetMap := make(map[string][]string)
	for _, edge := range edges {
		source := gjson.Get(edge.String(), `source`).String()
		target := gjson.Get(edge.String(), `target`).String()
		sourceToTargetMap[source] = append(sourceToTargetMap[source], target)
	}
	return sourceToTargetMap
}

// 搜索得到变量名称
func GetVariableName(str string) string {
	pattern := `{{loopData[@.]{1}(.*)}}`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(str)

	if len(match) > 1 {
		return match[1]
	} else {
		return ""
	}
}

// 替换所有的变量
func ReplaceAllVariable(str string, variables *simplejson.Json) string {
	str = strings.ReplaceAll(str, "$push:", "")
	for k := range variables.MustMap() {
		if strings.Contains(str, k) {
			valueTmp, _ := variables.Get(k).MarshalJSON()
			value := string(valueTmp)
			str = strings.ReplaceAll(str, k, value)
		}
	}
	return str
}

func CssToXpath(css string) string {
	css = strings.TrimSpace(css)
	parts := strings.Split(css, ">")
	xpath := "//"

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "#") {
			xpath += "*[@id='" + part[1:] + "']"
		} else if strings.HasPrefix(part, ".") {
			xpath += "*[contains(concat(' ', normalize-space(@class), ' '), ' " + part[1:] + " ')]"
		} else {
			if strings.Contains(part, ".") {
				tagAndClass := strings.Split(part, ".")
				xpath += tagAndClass[0] + "[contains(concat(' ', normalize-space(@class), ' '), ' " + tagAndClass[1] + " ')]"
			} else {
				xpath += part
			}
		}
		xpath += "/"
	}

	return xpath[:len(xpath)-1]
}

// Remove extra symbols
func RemoveExtraTextContent(text string) string {
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\t", "", -1)
	return text
}
