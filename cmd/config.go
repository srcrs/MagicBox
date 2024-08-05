package cmd

import (
	"MagicBox/utils"
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "about config",
	}

	configInit = &cobra.Command{
		Use:   "init",
		Short: "init task",
	}

	chromedpCtx context.Context

	exampleMap = map[string]string{
		"hostloc_get_integral": "examples/hostloc_get_integral.json",
		"jd_apply_refund":      "examples/jd_apply_refund.json",
		"v2ex_sign":            "examples/v2ex_sign.json",
		"wxread_task":          "examples/wxread_task.json",
		"bilibili_task":        "examples/bilibili_task.json",
	}

	UserInput = &InputParams{}
)

type InputParams struct {
	UserName  string
	PassWord  string
	BarkUrl   string
	Cron      string
	Cookies   string
	IPAddress string
}

func init() {
	configCmd.PersistentFlags().StringVarP(&UserInput.UserName, "username", "", "", "config login username")
	configCmd.PersistentFlags().StringVarP(&UserInput.PassWord, "password", "", "", "config login password")
	configCmd.PersistentFlags().StringVarP(&UserInput.Cron, "cron", "", "", "scheduled execution time")
	configCmd.PersistentFlags().StringVarP(&UserInput.BarkUrl, "barkUrl", "", "", "config notify bark")
	configCmd.PersistentFlags().StringVarP(&UserInput.Cookies, "cookies", "", "", "config cookies")
	configCmd.PersistentFlags().StringVarP(&UserInput.IPAddress, "ip", "", "", "remote chrome ip")
	configCmd.AddCommand(configInit)
	configInit.AddCommand(configByHostLocGetIntegral)
	configInit.AddCommand(configByJdApplyRefund)
	configInit.AddCommand(configByV2exSign)
	configInit.AddCommand(configByWxReadTask)
	configInit.AddCommand(configByBilibiliTask)
}

var configByBilibiliTask = &cobra.Command{
	Use:   "bilibili_task",
	Short: "bilibili everyday task, login, watch video",
	Run: func(cmd *cobra.Command, args []string) {
		config := utils.ChromeConfig{
			Headless: true,
			Timeout:  180,
		}
		taskName := "bilibili_task"
		indexUrl := "https://www.bilibili.com/"
		diffUrl := "https://account.bilibili.com/account/home"
		initConfigFile(taskName, config, indexUrl, diffUrl, "", []string{})
	},
}

var configByV2exSign = &cobra.Command{
	Use:   "v2ex_sign",
	Short: "v2ex sign",
	Run: func(cmd *cobra.Command, args []string) {
		config := utils.ChromeConfig{
			Headless: true,
			Timeout:  180,
		}
		taskName := "v2ex_sign"
		indexUrl := "https://www.v2ex.com"
		query := `//*[@class="tools"]`
		keyWords := []string{`注册`, `登录`}
		initConfigFile(taskName, config, indexUrl, "", query, keyWords)
	},
}

var configByWxReadTask = &cobra.Command{
	Use:   "wxread_task",
	Short: "wxread daily challenge task",
	Run: func(cmd *cobra.Command, args []string) {
		config := utils.ChromeConfig{
			Headless: true,
			Timeout:  180,
		}
		taskName := "wxread_task"
		indexUrl := "https://weread.qq.com"
		query := `//*[@class="navBar_border"]`
		keyWords := []string{`登录`}
		initConfigFile(taskName, config, indexUrl, "", query, keyWords)

	},
}

var configByHostLocGetIntegral = &cobra.Command{
	Use:   "hostloc_get_integral",
	Short: "hostloc get integral",
	Run: func(cmd *cobra.Command, args []string) {
		taskName := "hostloc_get_integral"
		fileContent, err := os.ReadFile(exampleMap[taskName])
		if err != nil {
			utils.GLOBAL_LOGGER.Error("hostloc file read error: " + err.Error())
			return
		}
		if UserInput.UserName == "" || UserInput.PassWord == "" {
			utils.GLOBAL_LOGGER.Warn("please enter username or password")
			return
		}

		newFileContent := userInputReplaceToFile(string(fileContent))
		utils.WriteToFile(exampleMap[taskName], []byte(newFileContent))
	},
}

var configByJdApplyRefund = &cobra.Command{
	Use:   "jd_apply_refund",
	Short: "auto apply refund",
	Run: func(cmd *cobra.Command, args []string) {
		config := utils.ChromeConfig{
			Headless: true,
			Timeout:  180,
		}
		taskName := "jd_apply_refund"
		indexUrl := "https://www.jd.com"
		diffUrl := "https://pcsitepp-fm.jd.com/"
		initConfigFile(taskName, config, indexUrl, diffUrl, "", []string{})
	},
}

func getCookies(config utils.ChromeConfig, indexUrl, diffUrl, query string, keyWords []string) string {
	utils.GLOBAL_LOGGER.Info("It will close in " + strconv.Itoa(config.Timeout-30) + " seconds")
	chromedpCtx, cancel, err := utils.NewChromedpContext(config)
	defer cancel()
	if err != nil {
		utils.GLOBAL_LOGGER.Error("init chrome error: " + err.Error())
		return ""
	}
	loginFlag := false
	jsonResult := ""
	if err := chromedp.Run(
		chromedpCtx,
		chromedp.Navigate("http://localhost:9222/json"),
		chromedp.InnerHTML("/html/body/pre", &jsonResult),
		chromedp.Navigate(indexUrl),
	); err != nil {
		utils.GLOBAL_LOGGER.Error("err: " + err.Error())
	}
	visitUrl := "http://localhost:9222" + gjson.Get(jsonResult, `0.devtoolsFrontendUrl`).String()
	if UserInput.IPAddress != "" {
		visitUrl = strings.ReplaceAll(visitUrl, "localhost", UserInput.IPAddress)
	}
	utils.GLOBAL_LOGGER.Info("please visit url: " + visitUrl)
	if err := chromedp.Run(
		chromedpCtx,
		chromedp.Sleep(time.Duration(config.Timeout-30)*time.Second),
	); err != nil {
		utils.GLOBAL_LOGGER.Error("sleep error: " + err.Error())
	}
	if diffUrl != "" {
		if err := chromedp.Run(
			chromedpCtx,
			utils.CheckLoginStatusBySkipLogin(&loginFlag, diffUrl),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("check login error: " + err.Error())
		}
	} else {
		if err := chromedp.Run(
			chromedpCtx,
			utils.CheckLoginStatusByKeyWords(&loginFlag, indexUrl, query, keyWords...),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("check login error: " + err.Error())
		}
	}
	utils.GLOBAL_LOGGER.Info("login result", zap.Any("status", loginFlag))
	if !loginFlag {
		cancel()
		return ""
	}
	cookiesResult := ""
	if err := chromedp.Run(
		chromedpCtx,
		utils.PrintCookies(&cookiesResult),
	); err != nil {
		utils.GLOBAL_LOGGER.Error("print cookie error: ", zap.Error(err))
	}
	return cookiesResult
}

func userInputReplaceToFile(fileContent string) string {
	var err error
	if UserInput.UserName != "" {
		nodeOld := gjson.Get(fileContent, `drawflow.nodes.#(label=="forms")#|#(data.value=="username")`).String()
		nodeNew, err := sjson.Set(nodeOld, `data.value`, UserInput.UserName)
		fileContent = strings.ReplaceAll(fileContent, nodeOld, nodeNew)
		if err != nil {
			utils.GLOBAL_LOGGER.Error("sjson set username file error: " + err.Error())
		}
	}
	if UserInput.PassWord != "" {
		nodeOld := gjson.Get(fileContent, `drawflow.nodes.#(label=="forms")#|#(data.value=="password")`).String()
		nodeNew, err := sjson.Set(nodeOld, `data.value`, UserInput.PassWord)
		fileContent = strings.ReplaceAll(fileContent, nodeOld, nodeNew)
		if err != nil {
			utils.GLOBAL_LOGGER.Error("sjson set password file error: " + err.Error())
		}
	}
	if UserInput.Cron != "" {
		fileContent, err = sjson.Set(fileContent, `drawflow.nodes.#(label=="trigger").data.triggers.#(type="cron-job").data.expression`, UserInput.Cron)
		if err != nil {
			utils.GLOBAL_LOGGER.Error("sjson set cron file error: " + err.Error())
		}
	}
	if UserInput.BarkUrl != "" {
		fileContent, err = sjson.Set(fileContent, `drawflow.nodes.#(label=="webhook")#.data.url`, UserInput.BarkUrl)
		if err != nil {
			utils.GLOBAL_LOGGER.Error("sjson set barkUrl file error: " + err.Error())
		}
	}
	if UserInput.Cookies != "" {
		fileContent, err = sjson.Set(fileContent, `drawflow.nodes.#(label=="insert-data")#.data.dataList.#(name="cookies").value`, UserInput.Cookies)
		if err != nil {
			utils.GLOBAL_LOGGER.Error("sjson set Cookies file error: " + err.Error())
		}
	}
	return fileContent
}

func initConfigFile(taskName string, config utils.ChromeConfig, indexUrl, diffUrl, query string, keyWords []string) {
	cookiesResult := getCookies(config, indexUrl, diffUrl, query, keyWords)
	fileContent, err := os.ReadFile(exampleMap[taskName])
	if err != nil {
		utils.GLOBAL_LOGGER.Error("file read error: " + err.Error())
		return
	}

	if cookiesResult == "" {
		return
	}

	newFileContent, err := sjson.Set(string(fileContent), `drawflow.nodes.#(label=="insert-data")#.data.dataList.#(name="cookies").value`, cookiesResult)
	if err != nil {
		utils.GLOBAL_LOGGER.Error("sjson set cookie file error: " + err.Error())
	}
	newFileContent = userInputReplaceToFile(newFileContent)
	utils.WriteToFile(exampleMap[taskName], []byte(newFileContent))
}
