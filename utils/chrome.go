package utils

import (
	"context"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

type ChromeConfig struct {
	//设备类型
	Device    string
	UserAgent string
	//是否开启无头模式
	Headless bool
	//屏幕大小
	WindowSizeWidth  int
	WindowSizeHeight int
	//本次执行超时时间 int
	Timeout int
}

func GetChromeConfigOpts(config ChromeConfig) []func(*chromedp.ExecAllocator) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		//不检查默认浏览器
		chromedp.NoDefaultBrowserCheck,
		//是否不开启图像界面
		chromedp.Flag("headless", config.Headless),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		//忽略证书错误
		chromedp.Flag("ignore-certificate-errors", true),
		//禁用网络安全标志
		chromedp.Flag("disable-web-security", true),
		chromedp.WindowSize(1920, 1080),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("enable-webgl", true),
		chromedp.Flag("ignore-gpu-blocklist", true),
		chromedp.Flag("disable-gpu-driver-bug-workarounds", true),
		chromedp.Flag("disable-gpu-vsync", true),
		chromedp.Flag("disable-gpu-sandbox", true),
		chromedp.Flag("disable-accelerated-video-decode", true),
		chromedp.Flag("disable-accelerated-video-encode", true),
		chromedp.Flag("disable-gpu-memory-buffer-video-frames", true),
		chromedp.Flag("disable-gpu-rasterization", true),
		chromedp.Flag("disable-2d-canvas-clip-aa", true),
		chromedp.Flag("disable-2d-canvas-image-chromium", true),
		chromedp.Flag("disable-gpu-compositing", true),
		chromedp.Flag("disable-gpu-shader-disk-cache", true),
		chromedp.Flag("disable-threaded-animation", true),
		chromedp.Flag("disable-threaded-scrolling", true),
		chromedp.Flag("disable-webgl-image-chromium", true),
		chromedp.Flag("disable-webgl", false),
		chromedp.Flag("use-gl", "desktop"),
		chromedp.Flag("enable-gpu-service-logging", true),
		chromedp.Flag("use-angle", "gl"),
		chromedp.Flag("use-cmd-decoder", "passthrough"),
		chromedp.Flag("use-gl", "desktop"),
		chromedp.Flag("enable-gpu-rasterization", true),
		chromedp.Flag("remote-debugging-port", "9222"),
		chromedp.Flag("remote-debugging-address", "0.0.0.0"),
		chromedp.Flag("remote-allow-origins", "*"),
		//设置网站不是首次运行
		chromedp.NoFirstRun,
		//设置代理
		// chromedp.ProxyServer(GetProxyUrl()),
		//设置UserAgent
		chromedp.UserAgent(GetUserAgent()),
	)
	return opts
}

// 获取cookies
// A=1; B=2
func GetCookiesJson2String(result *string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		// 1. 获取cookies
		cookies, err := network.GetCookies().Do(ctx)
		if err != nil {
			GLOBAL_LOGGER.Error("获取cookie失败", zap.Error(err))
			return err
		}
		// 输出 cookie
		for _, cookie := range cookies {
			*result += cookie.Name + "=" + cookie.Value + "; "
		}
		return nil
	}
}

// 加载cookie到chrome
func LoadCookies(cookiesData string) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		if cookiesData == "" {
			GLOBAL_LOGGER.Error("cookie为空", zap.Any("callid", ctx.Value("callid")))
			return nil
		}
		// 反序列化
		cookiesParams := network.SetCookiesParams{}
		if err = cookiesParams.UnmarshalJSON([]byte(cookiesData)); err != nil {
			return nil
		}
		// 设置cookies
		return network.SetCookies(cookiesParams.Cookies).Do(ctx)
	}
}

// 打印cookie
func PrintCookies(cookiesResult *string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		// cookies的获取对应是在devTools的network面板中
		// 1. 获取cookies
		cookies, err := network.GetCookies().Do(ctx)
		if err != nil {
			return err
		}
		// 2. 序列化
		cookiesData, err := network.GetCookiesReturns{Cookies: cookies}.MarshalJSON()
		if err != nil {
			return err
		}
		*cookiesResult = string(cookiesData)
		// 打印cookie
		GLOBAL_LOGGER.Info("cookie: \n" + string(cookiesData))
		return nil
	}
}

// 判断是否登录, 根据url是否跳转登录页面
func CheckLoginStatusBySkipLogin(loginFlag *bool, url string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 40*time.Second)
		defer cancel()
		var location string
		if err := chromedp.Run(
			ctx,
			chromedp.Navigate(url),
			chromedp.Sleep(5*time.Second),
			chromedp.Location(&location),
		); err != nil {
			GLOBAL_LOGGER.Error("检查登录失败", zap.Error(err), zap.Any("callid", ctx.Value("callid")))
			return err
		}
		if location != url {
			*loginFlag = false
		} else {
			*loginFlag = true
		}
		GLOBAL_LOGGER.Info("当前登录状态为：", zap.Bool("login", *loginFlag))
		return nil
	}
}

// 判断是否登录, 根据url是否跳转登录页面
func CheckLoginStatusByKeyWords(loginFlag *bool, url, query string, words ...string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		var text string
		err := chromedp.Run(
			ctx,
			chromedp.Navigate(url),
			chromedp.Sleep(5*time.Second),
			chromedp.WaitVisible(query),
			chromedp.Sleep(5*time.Second),
			chromedp.Text(query, &text),
		)
		if err != nil {
			GLOBAL_LOGGER.Error("check login error", zap.Error(err), zap.Any("callid", ctx.Value("callid")))
			return err
		}
		if text != "" {
			*loginFlag = false
			for _, word := range words {
				if !strings.Contains(text, word) {
					*loginFlag = true
					break
				}
			}
		}
		GLOBAL_LOGGER.Info("login result: ", zap.Bool("login", *loginFlag), zap.Any("callid", ctx.Value("callid")))
		return nil
	}
}

func GetAttributeText(node *cdp.Node) string {
	var text string
	for _, child := range node.Children {
		if child.NodeType == cdp.NodeTypeText {
			text += child.NodeValue
		} else {
			text += GetAttributeText(child)
		}
	}
	return text
}

func NewChromedpContext(config ChromeConfig) (context.Context, context.CancelFunc, error) {
	//chromedp初始参数
	opts := GetChromeConfigOpts(config)
	//创建一个上下文
	allCtx, cancel := chromedp.NewExecAllocator(
		context.Background(),
		opts...,
	)
	chromedpCtx, cancel := chromedp.NewContext(
		allCtx,
	)
	//创建超时时间
	chromedpCtx, cancel = context.WithTimeout(chromedpCtx, time.Duration(config.Timeout)*time.Second)

	return chromedpCtx, cancel, nil
}
