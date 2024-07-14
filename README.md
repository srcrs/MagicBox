<div align="center"> 
<h1 align="center">MagicBox</h1>
<img src="https://img.shields.io/github/issues/srcrs/MagicBox?color=green">
<img src="https://img.shields.io/github/stars/srcrs/MagicBox?color=yellow">
<img src="https://img.shields.io/github/forks/srcrs/MagicBox?color=orange">
<img src="https://img.shields.io/github/license/srcrs/MagicBox?color=ff69b4">
<img src="https://img.shields.io/github/search/srcrs/MagicBox/main?color=blue">
</div>

## 简述

`MagicBox`在今年迎来了升级，依托于`Automa`灵活的工作流配置，可以很方便的实现网站的自动化任务。`Automa`是一个浏览器控制插件，有着众多的组件，基本涵盖了日常的操作，只需要拖拉连线便可将打开网页、点击链接、获取元素的文本等等拼接成一个工作流话不，实现复杂任务的简化执行，但，如果想让其在服务端每日自动执行，不必依赖本地的插件环境，是否可行呢？

新版`MagicBox`的主要任务都会使用`Automa`来实现，实践的过程中是发现了一些问题的，例如登录态、通知、代码处理逻辑不一致等问题，对于迁移使用仍然会有一定的理解难度。

在最新的`2.2.2`版本中，我们新加了`cli`命令模式，内置了部分任务标准模版，只需根据命令引导，即可生成专属于自己的任务，极大的简化了使用`MagicBox`。

## 目录

- [简述](#简述)
- [目录](#目录)
- [项目目录说明](#项目目录说明)
- [已实现组件](#已实现组件)
- [内置支持任务](#内置支持任务)
- [环境要求](#环境要求)
  - [Linux](#linux)
    - [Docker](#docker)
- [使用方法](#使用方法)
  - [Docker部署](#docker部署)
    - [1.克隆仓库](#1克隆仓库)
    - [2.导入配置](#2导入配置)
      - [使用内置命令初始化配置](#使用内置命令初始化配置)
      - [自定义导入](#自定义导入)
    - [3.执行部署](#3执行部署)
- [使用示例](#使用示例)
  - [v2ex\_sign](#v2ex_sign)
  - [hostloc\_get\_integral](#hostloc_get_integral)
  - [jd\_apply\_refund](#jd_apply_refund)
  - [wxread\_task](#wxread_task)
- [开发贡献](#开发贡献)
  - [加载cookie](#加载cookie)
  - [定时执行](#定时执行)
  - [用户登录](#用户登录)
- [通知方式](#通知方式)
  - [Bark](#bark)

## 项目目录说明

项目地址：https://github.com/srcrs/MagicBox

```
MagicBox
├── Dockerfile
├── LICENSE
├── MagicBox.log
├── README.md
├── cmd
├── configs
├── docker-compose.yml
├── examples
├── go.mod
├── go.sum
├── install.sh
├── main.go
├── public
├── script.sh
├── utils
└── workerflow
```

- examples: 有示例配置文件通过cli命令可以重复使用
- configs: 用于放置需要执行的automa配置文件
- docker-compose.yml: docker-compose配置文件，实时获取最新的版本
- MagicBox.log: 工作流执行后的日志文件
- main.go: 工作流解析引擎执行入口
- utils、workerflow: 解析引擎相关核心代码逻辑

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

## 内置支持任务

- | 站点 | 说明 | 登录授权方式 | username | password | brakUrl | cron
-|-|-|-|-|-|-|-
hostloc_get_integral | https://hostloc.com/ | 每日访问空间刷积分 | 账号密码 | yes | yes | yes | yes
jd_apply_refund | https://www.jd.com/ | 京东自动申请价格保护 | cookie | no | no | yes | yes
v2ex_sign | https://v2ex.com/ | 每日签到 | cookie | no | no | yes | yes
wxread_task | https://weread.qq.com/ | 每日登录阅读，完成读书挑战 | cookie | no | no | yes | yes

## 环境要求

### Linux

#### Docker

docker环境安装参考[官方教程](https://docs.docker.com/engine/install/debian/)，一键把docker和docker-compose环境都安装好。

```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```

安装好后，示例docker版本信息。

```bash
$ docker --version
Docker version 24.0.4
$ docker compose version
Docker Compose version v2.19.1
```

## 使用方法

### Docker部署

#### 1.克隆仓库

```bash
git clone https://github.com/srcrs/MagicBox.git
```

克隆完后，进入到MagicBox文件夹中。

#### 2.导入配置

目前支持两种方式，通过内置任务初始化生成；导入自定义任务（可能存在部分节点未接入问题）。

##### 使用内置命令初始化配置

- 帮助命令

```bash
#查看目前支持的命令
$docker compose run --rm -p 9222:9222 server -h

#查看目前支持的config命令
$docker compose run --rm -p 9222:9222 server config -h

#查看目前支持的初始化任务
$docker compose run --rm -p 9222:9222 server config init -h
```

- 可传入参数

变量名 | 说明 | 使用示例
-|-|-
username | 登录用户名 | --username "xxxxxx"
password | 登录密码 | --password "xxxxxxx"
barkUrl | 通知 | --barkUrl "xxxxxxx"
cron | 定时执行 | --cron "12 12 * * *"

- 真实案例

初始化一个v2ex任务。

```bash
#1.选择初始化v2ex任务，设置定时执行时间和bark通知
$docker compose run --rm -p 9222:9222 server config init v2ex_sign --cron "12 12 * * *" --barkUrl "https://bark.xxx.com/xxxxxx"
It will close in 150 seconds
please visit url: http://localhost:9222/devtools/inspector.html?ws=localhost:9222/devtools/page/333BE3874077691C51A4279C7A4E8AB9

#2.将locahost替换为服务器ip后在浏览器访问，即可远程控制，在150秒内完成登录操作

#3.150秒后将会检查登录情况，会自动将配置文件添加到configs文件夹中
new config path: configs/v2ex_sign_7656fbd2-78dc-4cff-af18-e56c40b8e527.json
```

更多请参考：[使用示例](#使用示例)

##### 自定义导入

将automa编写好的配置从插件中导出，再导入至MagicBox中，适合有一定的使用经验。

#### 3.执行部署

```bash
#进入到MagicBox目录
docker compose up -d
```

## 使用示例

### v2ex_sign

```bash
$docker compose run --rm -p 9222:9222 server config init v2ex_sign --cron "12 12 * * *" --barkUrl "https://bark.xxx.com/xxxxxx"
```

### hostloc_get_integral

```bash
$docker compose run --rm -p 9222:9222 server config init hostloc_get_integral --cron "12 12 * * *" --barkUrl "https://bark.xxx.com/xxxxxx" --username "xxxxxxx" --password "yyyyyyy"
```

### jd_apply_refund

```bash
$docker compose run --rm -p 9222:9222 server config init jd_apply_refund --cron "12 12 * * *" --barkUrl "https://bark.xxx.com/xxxxxx"
```

### wxread_task

```bash
$docker compose run --rm -p 9222:9222 server config init wxread_task --cron "12 12 * * *" --barkUrl "https://bark.xxx.com/xxxxxx"
```

## 开发贡献

Automa是本地执行，在实际迁移使用时，需要考虑到登录态问题，定时任务、用户名和密码登录等，也有相应的使用规范。

### 加载cookie

使用`insert-data`组件，选择添加Variable变量，名称为`cookies`，一般要在页面打开之前将cookie加载进去。

### 定时执行

修改`Trigger`组件，添加一个`Cron job`，就可以设置cron定时任务了。MagicBox加载逻辑是，程序首次都会执行一次，后续是根据定时任务的设置执行。

### 用户登录

使用的组件是`Forms`，使用`Text field`填写内容，用户名是`username`，密码是`password`。

## 通知方式

automa的`HTTP request`可以实现接口的调用，正好可以满足通知的需求，但是通知的内容可能会很少，无法做到代码方式灵活。

### Bark

Bark 是一款`iOS`端的推送服务，通过部署一个`Server`服务端，可以通过HTTP接口来给`iOS`设备发送推送通知，代码开源: https://github.com/Finb/Bark

使用极其简单，只需要下载安装Bark软件即可，获取每个用户唯一key，按照下面格式替换，就得到了你唯一的推送通道了。一般将Bark的链接替换到组件中即可。

```bash
https://api.day.app/DnzTsd6qDWTdfs9xRGygFtasdnsRCL/
```

详细可参考：[Bark官方文档](https://bark.day.app/)