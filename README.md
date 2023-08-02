<div align="center"> 
<h1 align="center">MagicBox</h1>
<img src="https://img.shields.io/github/issues/srcrs/MagicBox?color=green">
<img src="https://img.shields.io/github/stars/srcrs/MagicBox?color=yellow">
<img src="https://img.shields.io/github/forks/srcrs/MagicBox?color=orange">
<img src="https://img.shields.io/github/license/srcrs/MagicBox?color=ff69b4">
<img src="https://img.shields.io/github/search/srcrs/MagicBox/main?color=blue">
</div>

## 简述

万事总有缘由，MagicBox也是一样，愿如百宝箱，简化世间万物，简而言之就是代替重复性任务。目前实现了doduo日常用的频次较高网站签到，接下来将继续接入更多的签到，代替更多的重复性任务。

关于安全性，可能没有比这更高的了，底层直接操作原生chrome浏览器，模拟用户行为，每次使用完自动更新cookie（理论达到续期目的），但同时会对机器的性能要求更高，启动chrome要占用更多的资源。

说明：若目前有需要使用贴吧签到功能，可以使用免费授权码（填写到default
.yml的token变量），目前稳定性测试中，到期时间`2023-08-24 22:28:11`

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTI4ODcyOTEsIm1heF90YXNrcyI6OTksIm1heF91c2VycyI6NTAwfQ.-bHXra6O5177hv0SDL-a2cil8iKevL_fvTRpPJIjvgc
```

- [简述](#简述)
- [项目目录解答](#项目目录解答)
  - [configs篇](#configs篇)
  - [本地获取cookie](#本地获取cookie)
  - [docker-compose.yml](#docker-composeyml)
  - [任务执行日志](#任务执行日志)
- [环境说明](#环境说明)
- [食用方法](#食用方法)
  - [自有机器Docker部署](#自有机器docker部署)
  - [通用获取cookie](#通用获取cookie)
- [支持任务列表](#支持任务列表)
  - [吾爱破解](#吾爱破解)
    - [目前已实现功能](#目前已实现功能)
    - [cookie获取方法](#cookie获取方法)
    - [配置文件示例](#配置文件示例)
  - [哔哩哔哩](#哔哩哔哩)
    - [目前已实现功能](#目前已实现功能-1)
    - [cookie获取方法](#cookie获取方法-1)
    - [配置文件示例](#配置文件示例-1)
  - [全球主机交流论坛](#全球主机交流论坛)
    - [目前已实现功能](#目前已实现功能-2)
    - [准备用户名/密码](#准备用户名密码)
    - [配置文件示例](#配置文件示例-2)
  - [v2ex论坛](#v2ex论坛)
    - [目前已实现功能](#目前已实现功能-3)
    - [cookie获取方法](#cookie获取方法-2)
    - [配置文件示例](#配置文件示例-3)
  - [贴吧签到（受限制）](#贴吧签到受限制)
    - [目前已实现功能](#目前已实现功能-4)
    - [cookie获取方法](#cookie获取方法-3)
    - [配置文件示例](#配置文件示例-4)
  - [有道云签到](#有道云签到)
    - [目前已实现功能](#目前已实现功能-5)
    - [cookie获取方法](#cookie获取方法-4)
    - [配置文件示例](#配置文件示例-5)
- [通知方式](#通知方式)
  - [Bark](#bark)
  - [Telegram](#telegram)
  - [企业微信应用通知](#企业微信应用通知)

## 项目目录解答

项目地址：https://github.com/srcrs/MagicBox

```
MagicBox
├── configs
│   ├── 52pojie.yml
│   ├── bilibili.yml
│   ├── default.yml
│   ├── hostloc.yml
│   └── v2ex.yml
├── docker-compose.yml
├── LICENSE
├── MagicBox_amd64_darwin
├── MagicBox_amd64_linux
├── MagicBox_amd64_win.exe
├── MagicBox.log
└── README.md
```

### configs篇

将项目克隆到本地后，会得到类似目录，configs目录中都是配置信息，default.yml配置项目的默认配置，优先级最低；此外不同站点的任务执行配置也是拆分开来的，这里主要考虑到是便于管理，部分任务执行所需信息内容过多，如cookie等，再者支持多用户的使用，内容多了也不便于配置。

default.yml支持的配置

```yml
notify: '通知'
token: '开启高级权限'
```

站点的配置文件，如果不需要此任务直接删除即可，接下来看一个配置示例。

```yml
bilibili:
  users:
    doduo:
      cookie: ''
      cron: '0 41 8,16 * * *'
  task:
    watchVideo: true
    thowCoin: 5
    shareVideo: true
    notify: ''
    multiThread: false
```

`bilibili`代表目标任务，下面有两个大层级`users`和`task`，`users`下面可以多个用户复制`doduo`层级即可，`task`代表这些用户需要做哪些任务、推送地址、是否支持并发执行。`task`下的配置，支持在用户层级配置，优先级最高。

```yml
bilibili:
  users:
    doduo:
      cookie: ''
      cron: '0 41 8,16 * * *'
    doduo2:
      cookie: ''
      cron: '0 41 8,16 * * *'
      thowCoin: 1
      notify: 'yyyyyy'
  task:
    watchVideo: true
    thowCoin: 5
    shareVideo: true
    notify: 'xxxxxx'
    multiThread: false
```

示例表示：需要执行哔哩哔哩任务，默认需要执行观看视频、每日投币5个、分享视频，并设置了推送通知，有两个用户要执行doduo和doduo2，其中doduo执行默认任务，doduo2在此基础上自定义了每日投币1个，以及推送通知。

### 本地获取cookie

|二进制文件|支持平台
-|-
MagicBox_amd64_darwin|mac平台
MagicBox_amd64_linux|linux平台
MagicBox_amd64_win.exe|win平台

用于在各个平台本地获取cookie，在任务详细篇再一一介绍。

### docker-compose.yml

docker-compose配置文件，便于重复使用。

### 任务执行日志

存储在MagicBox.log中。

## 环境说明

- 程序底层依赖chrome浏览器，因此若涉及到本地操作需要有该环境

- 二进制文件包类别

文件名 | 对应平台
-|-
MagicBox_amd64_darwin | mac
MagicBox_amd64_linux | linux
MagicBox_amd64_win.exe | win

若是在win平台获取cookie，和其他两个有些不同。在与可执行文件（MagicBox_amd64_win.exe）同目录下，创建exec.bat，内容如下

```bat
start MagicBox_amd64_win.exe bilibili login
```

这个是bilibili的示例，无win设备，如有更好的办法欢迎推荐。

## 食用方法

### 自有机器Docker部署

- 1.克隆仓库

```bash
git clone git@github.com:srcrs/MagicBox.git
```

- 2.在configs目录下填写对应任务的配置文件，参考各个任务详细介绍

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

### 通用获取cookie

支持v0.1.3版本以上

- mac

```bash
./MagicBox_amd64_darwin chrome login
```

- win

```bash
./MagicBox_amd64_win.exe chrome login
```

- linux

```bash
./MagicBox_amd64_linux chrome login
```

## 支持任务列表

### 吾爱破解

官方站点：https://www.52pojie.cn/

#### 目前已实现功能

- 每日签到

#### cookie获取方法

吾爱登录有人机验证，使用cookie较为简单，需要在本地电脑获取，确保有chrome浏览器。

- mac示例

根据提示快速进行登录以获取cookie

```bash
./MagicBox_amd64_darwin 52pojie login
```

- docker示例

```bash
docker compose run server 52pojie login
```

#### 配置文件示例

文件名：52pojie.yml，示例每天8:41、16:41各执行一次。

```yml
52pojie:
  users:
    doduo:
      cookie: 'xxxxxxxx'
      cron: '0 41 8,16 * * *'
  task:
    multiThread: false
    notify: 'yyyyyyyyy'
    checkIn: true
```

变量配置说明

|变量名|说明|
-|-
cookie | 52pojie的cookie信息
checkIn | 是否签到
cron | cron执行任务
notify | 通知
multiThread | 是否支持并发，默认填写false即可

### 哔哩哔哩

官方站点：https://www.bilibili.com/

#### 目前已实现功能

- 每日登录

- 每日投币

- 每日观看视频

- 每日分享视频

- 每日直播签到

视频选取策略：优先取用户动态列表和综合热门视频。

#### cookie获取方法

- mac示例

目前支持扫码获取cookie

```
./MagicBox_amd64_darwin bilibili login
```

- docker示例

```bash
docker compose run server bilibili login
```

#### 配置文件示例

配置文件名：bilibili.yml

```yml
bilibili:
  users:
    doduo:
      cookie: 'xxxxxxxxxx'
      cron: '0 41 8,16 * * *'
  task:
    watchVideo: true
    thowCoin: 5
    shareVideo: true
    notify: ''
    multiThread: false
```

变量配置说明

|变量名|说明|
-|-
cookie | bilibili的cookie信息
cron | cron执行任务
watchVideo | 是否每天观看视频
thowCoin | 每天投币数量
shareVideo | 是否分享视频
notify | 通知
multiThread | 是否支持并发，默认填写false即可

### 全球主机交流论坛

官方站点：https://hostloc.com/

#### 目前已实现功能

- 每日访问空间10次获得积分

访问策略：选取前两页出现的用户，依次进行访问

- 支持跳过访问某些用户空间

#### 准备用户名/密码

可以使用用户名和密码进行登录

#### 配置文件示例

```yml
hostloc:
  users:
    doduo:
      username: ''
      password: ''
      cron: '0 41 8,16 * * *'
  task:
    accessSpace: 10
    notify: ''
    filersUser: "xxxxx, yyyyyy, zzzzzzzz"
    multiThread: false
```

变量配置说明

|变量名|说明|
-|-
username | 用户名
password | 密码
cron | cron执行任务
accessSpace | 每次执行访问空间次数
filersUser | 过滤掉某些用户空间访问（避免每天重复访问）
notify | 通知
multiThread | 是否支持并发，默认填写false即可

### v2ex论坛

官方站点：https://v2ex.com/

#### 目前已实现功能

- 每日签到

#### cookie获取方法

v2ex登录有图片验证码缓解，将base64打印到控制台，用在线网站可以查看图片信息，例如：https://tool.jisuapi.com/base642pic.html，然后根据提示输入验证码、用户名、密码。

- mac平台

```bash
./MagicBox_amd64_darwin v2ex login
```

- docker平台

```bash
docker compose run server v2ex login
```

#### 配置文件示例

```yml
v2ex:
  users:
    doduo:
      cookie: ''
      cron: '0 41 8,16 * * *'
  task:
    checkIn: true
    multiThread: false
    notify: ''
```

变量配置说明

|变量名|说明|
-|-
cookie | v2ex的cookie信息
checkIn | 是否签到
cron | cron执行任务
notify | 通知
multiThread | 是否支持并发，默认填写false即可

### 贴吧签到（受限制）

官方站点：https://tieba.baidu.com/

#### 目前已实现功能

- 每日签到（不限制数量）

  [签到经验获取规则](https://tieba.baidu.com/f/like/level?kw=&ie=utf-8&lv_t=lv_nav_who)，本程序模拟pc端签到。此功能需要填写授权码才可进行使用。

#### cookie获取方法

由于百度时常会遇到安全验证，因此建议本地获取cookie后进行使用。

- mac平台

```bash
./MagicBox_amd64_darwin tieba login
```

#### 配置文件示例

```yml
tieba:
  task:
    multiThread: false
    notify: ''
  users:
    doduo:
      cookie: ''
      cron: '0 41 8,16 * * *'
```

变量配置说明

|变量名|说明|
-|-
cookie | 密码
cron | cron执行任务
notify | 通知
multiThread | 是否支持并发，默认填写false即可

### 有道云签到

#### 目前已实现功能

- 每日签到（获得空间容量）

#### cookie获取方法

使用通用的方法获取，手动进行登录，90秒后将会吧cookie打印在控制台

- mac平台

```bash
./MagicBox_amd64_darwin chrome login
```

#### 配置文件示例

```yml
tieba:
  task:
    multiThread: false
    notify: ''
    checkIn: true
  users:
    doduo:
      cookie: ''
      cron: '0 41 8,16 * * *'
```

变量配置说明

|变量名|说明|
-|-
cookie | 密码
cron | cron执行任务
notify | 通知
multiThread | 是否支持并发，默认填写false即可
checkIn | 是否执行签到

## 通知方式

将通知链接填写到notify变量中即可，推送通道这里理论上都是可以支持的，由于个人精力有限，目前只对接了一部分，未来会随着小伙伴们的喜好进行迭代升级。

### Bark

Bark 是一款`iOS`端的推送服务，通过部署一个`Server`服务端，可以通过HTTP接口来给`iOS`设备发送推送通知，代码开源: https://github.com/Finb/Bark

使用极其简单，只需要下载安装Bark软件即可，获取每个用户唯一key，按照下面格式替换，就得到了你唯一的推送通道了。

```bash
#DnzTsd6qDWTdfs9xRGygFtasdnsRCL 需要替换的key
bark://api.day.app/DnzTsd6qDWTdfs9xRGygFtasdnsRCL/
```

详细可参考：[Bark官方文档](https://bark.day.app/)

### Telegram

这款软件在国内需要些特殊方式才可正常使用，若以前未使用过建议就不要折腾了。

- 打开bot的生成链接 `https://t.me/botfather`

- 点击 `/newbot - create a new bot` 生成新一个的bot，系统会让你给它取一个名字，这里可以随便输入，反正方便记就可以了。

- 系统会再次让你取一个名字并输入，现在的名字必须以bot结尾，不可以和其他任何bot重名。

- 一旦你输入的bot名字可用，系统会生成一个token给你，类似于`1729581149:BHGYVVjEHsaNjsnT8eQpWyshwr2o4PqU7u8`，请务必保存好此token并且不泄露，这是唯一的用户凭证。

- 此时，你拥有了一个bot，但是还无法使用，因为你知道它，它不知道你。这时候打开 `https://t.me/iamthebot` 这个链接，注意iamthebot为你刚才新建的bot的名字！点击/start进入对话框，发送 `@userinfobot` 后并点击它。`userinfobot`的对话中，点击或者输入`/start`，你将获取一个`Id/chat_id`，具体表现为一串数字，比如`387980691`。

执行完这些操作之后，这个机器人便已经搭建好了，示例推送格式如下：

```bash
tgram://1729581149:BHGYVVjEHsaNjsnT8eQpWyshwr2o4PqU7u8/387980691/
```

### 企业微信应用通知

示例推送格式如下：

```bash
qywx://sdadas:1023837:5DryRSrtLiasdsddsfsfZqPXaQaIajbSfO1trY
```

对上面格式解释

值 | 说明
-|-
sdadas | corpid
1023837 | agentid
5DryRSrtLiasdsddsfsfZqPXaQaIajbSfO1trY | corpsecret

corpid 和 corpsecret 主要用于获取access_token，可以参考微信文档[获取access_token](https://developer.work.weixin.qq.com/document/path/91039)填写内容，agentid是为了能发消息，可以看文档[发送应用消息](https://developer.work.weixin.qq.com/document/path/90236)。