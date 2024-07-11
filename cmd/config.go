package cmd

import (
	"MagicBox/utils"
	"context"
	"os"
	"strconv"
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
		Run: func(cmd *cobra.Command, args []string) {
			utils.GLOBAL_LOGGER.Info("hello")

		},
	}

	configInit = &cobra.Command{
		Use: "init",
	}

	chromedpCtx context.Context

	exampleMap = map[string]string{
		"hostloc_get_integral": "examples/hostloc_get_integral.json",
		"jd_apply_refund":      "examples/jd_apply_refund.json",
		"v2ex_sign":            "examples/v2ex_sign.json",
		"wxread_task":          "examples/wxread_task.json",
	}
)

func init() {
	configCmd.AddCommand(configInit)
	configInit.AddCommand(configByHostLocGetIntegral)
	configInit.AddCommand(configByJdApplyRefund)
	configInit.AddCommand(configByV2exSign)
	configInit.AddCommand(configByWxReadTask)
}

var configByV2exSign = &cobra.Command{
	Use:   "v2ex_sign",
	Short: "v2ex sign",
	Run: func(cmd *cobra.Command, args []string) {
		config := utils.ChromeConfig{
			Headless: true,
			Timeout:  150,
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
			Timeout:  150,
		}
		taskName := "wxread_task"
		indexUrl := "https://weread.qq.com"
		query := `//*[@class="navBar_border"]`
		keyWords := []string{`登录`}
		initConfigFile(taskName, config, indexUrl, "", query, keyWords)

	},
}

var configByJdApplyRefund = &cobra.Command{
	Use:   "jd_apply_refund",
	Short: "auto apply refund",
	Run: func(cmd *cobra.Command, args []string) {
		config := utils.ChromeConfig{
			Headless: true,
			Timeout:  150,
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

func initConfigFile(taskName string, config utils.ChromeConfig, indexUrl, diffUrl, query string, keyWords []string) {
	cookiesResult := getCookies(config, indexUrl, diffUrl, query, keyWords)
	fileContent, err := os.ReadFile(exampleMap[taskName])
	if err != nil {
		utils.GLOBAL_LOGGER.Error("file read error: " + err.Error())
	}

	newFileContent, err := sjson.Set(string(fileContent), `drawflow.nodes.#(label=="insert-data")#.data.dataList.#(name="cookies").value`, cookiesResult)
	if err != nil {
		utils.GLOBAL_LOGGER.Error("sjson ser file error: " + err.Error())
	}
	utils.WriteToFile(exampleMap[taskName], []byte(newFileContent))
}

var configByHostLocGetIntegral = &cobra.Command{
	Use: "hostloc",
	Run: func(cmd *cobra.Command, args []string) {
		utils.GLOBAL_LOGGER.Info("It will close in 120 seconds")

		config := utils.ChromeConfig{
			Headless: true,
			Timeout:  12000,
		}
		chromedpCtx, cancel, err := utils.NewChromedpContext(config)
		defer cancel()
		if err != nil {
			utils.GLOBAL_LOGGER.Error("init chrome error: " + err.Error())
			return
		}
		loginFlag := false

		indexUrl := "https://hostloc.com"
		query := `//*[@id="hd"]/div/div[1]`
		keyWords := []string{`自动登录`, `找回密码`, `注册`}
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
		utils.GLOBAL_LOGGER.Info("please visit url: " + visitUrl)
		if err := chromedp.Run(
			chromedpCtx,
			chromedp.Sleep(time.Duration(config.Timeout-90)*time.Second),
			utils.CheckLoginStatusByKeyWords(&loginFlag, indexUrl, query, keyWords...),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("check login error: " + err.Error())
		}
		utils.GLOBAL_LOGGER.Info("login result", zap.Any("status", loginFlag))
		if !loginFlag {
			cancel()
			return
		}
		cookiesResult := ""
		if err := chromedp.Run(
			chromedpCtx,
			utils.PrintCookies(&cookiesResult),
		); err != nil {
			utils.GLOBAL_LOGGER.Error("print cookie error: ", zap.Error(err))
		}
		fileContent, err := os.ReadFile(exampleMap["HostLocGetIntegral"])
		if err != nil {
			utils.GLOBAL_LOGGER.Error("file read error: " + err.Error())
		}

		newFileContent, err := sjson.Set(string(fileContent), `drawflow.nodes.#(label=="insert-data")#.data.dataList.#(name="cookies").value`, cookiesResult)
		if err != nil {
			utils.GLOBAL_LOGGER.Error("sjson ser file error: " + err.Error())
		}
		utils.WriteToFile(exampleMap["HostLocGetIntegral"], []byte(newFileContent))

		cancel()

	},
}
