<div align="center"> 
<h1 align="center">MagicBox</h1>
<img src="https://img.shields.io/github/issues/srcrs/MagicBox?color=green">
<img src="https://img.shields.io/github/stars/srcrs/MagicBox?color=yellow">
<img src="https://img.shields.io/github/forks/srcrs/MagicBox?color=orange">
<img src="https://img.shields.io/github/license/srcrs/MagicBox?color=ff69b4">
<img src="https://img.shields.io/github/search/srcrs/MagicBox/main?color=blue">
</div>

## 简述

[v1](https://github.com/srcrs/MagicBox/tree/v1)版本是使用代码来操作浏览器，写过几个自动化任务之后，发现流程极其相似，将浏览器操作颗粒化之后，能否使用工作流来实现？[Automa](https://github.com/AutomaApp/automa)便是最佳的选择，但局限于它是一个浏览器插件，无法在浏览器headless模式导入编写好的工作流，遂做了一个golang版本的工作流解析器，将Automa工作流导入到该项目中便可自动执行，以期平替其在本地化的操作，这便是v2版本。目前只实现了一部分操作，正在逐渐开发完善中。

## 已实现组件

- conditions：条件判断
- event-click：点击
- get-text：获取文本
- insert-data：插入变量
- loop-data：循环获取数据
- new-tab：打开网页
- webhook：调用接口
- tab-url：获取当前页面url
- element-scroll：滚动页面到屏幕最下面
- delay：流程sleep
- loop-elements：循环遍历页面元素
- forms：设置form表单填写内容
- reload-tab：刷新当前页面
- close-tab：关闭当前页面
- link：获取网页中链接打开页面
- active-tab：回到活动tab页中

## 目录

- [简述](#简述)
- [已实现组件](#已实现组件)
- [目录](#目录)
- [本地获取cookie](#本地获取cookie)
- [环境说明](#环境说明)
- [食用方法](#食用方法)
  - [Docker部署](#docker部署)
- [任务示例](#任务示例)
  - [v2ex论坛签到](#v2ex论坛签到)
  - [百度热搜自动推送](#百度热搜自动推送)
  - [京东自动申请价保](#京东自动申请价保)
  - [hostloc获取积分](#hostloc获取积分)
  - [微信读书完成每日阅读任务](#微信读书完成每日阅读任务)

## 本地获取cookie

涉及到登录问题，通常使用cookie来解决，推荐使用插件[cookie-editor](https://cookie-editor.com/)来获取，导出为json。

## 环境说明

- 程序底层依赖chrome浏览器，需要有该环境

- go 1.18

- docker

## 食用方法

### Docker部署

- 1.克隆仓库

```bash
git clone https://github.com/srcrs/MagicBox.git
```

- 2.在configs目录下导入对应任务的配置文件

- 3.安装docker环境

docker环境安装参考[官方教程](https://docs.docker.com/engine/install/debian/)，一键把docker和docker-compose环境都安装好

```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```

- 3.执行部署

示例docker版本信息

```bash
$ docker --version
Docker version 24.0.4
$ docker compose version
Docker Compose version v2.19.1
```

部署

```bash
docker compose up -d
```

## 任务示例

### v2ex论坛签到

官方站点：https://v2ex.com/

```
./configs/v2ex_sign.json
```

需要补充cookie，以及通知。

![](public/img/v2ex_sign.png)

### 百度热搜自动推送

官方站点：https://top.baidu.com/board?tab=realtime

```
./configs/post_notify.json
```

### 京东自动申请价保

官方站点：https://www.jd.com/

```
./configs/jd_sign.json
```

![](public/img/jd_sign.png)

### hostloc获取积分

官方站点：https://hostloc.com/

```
./configs/hostloc_sign.json
```

![](public/img/hostloc_sign.png)

### 微信读书完成每日阅读任务

官方站点：https://weread.qq.com/

需要填写cookie，以及需要阅读的书籍。

```
./configs/wxread_sign.json
```

![](public/img/wxread_task.png)