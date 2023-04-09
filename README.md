# 功能介绍

- 每日登录

- 每日观看视频

- 每日投币

- 每日分享视频

# 使用方法

目前程序已经打包成Docker镜像，可以很方便在Docker环境中使用。初次执行需要进行扫码登录，用户信息完全保存在本地。

**注意**：需要具备Docker、Docker-Compose环境

1. 克隆仓库

```bash
git clone https://github.com/asksowhat/bilibili-task-docker.git
```

2. 初次运行

进入`bilibili-task-docker`目录，执行下面命令，第一次需要扫码，后面会自动保存并更新cookie信息，理论上是不会失效的。

```bash
bash exec.sh
```

3. 自定义执行计划任务

首先修改`exec.sh`执行脚本，将 `/xxxx/yyyy/bilibili-task-docker`替换为你自己的真实目录

```
cd /xxxx/yyyy/bilibili-task-docker

docker-compose down

docker-compose up;docker-compose down
```

使用crontab，达到每日自动执行的目的，下面是一种示例格式

每小时执行一次：10:29、11:29、12:29......

```bash
29 */1 * * * /xxxx/yyyy/bilibili-task-docker/exec.sh > /xxxx/yyyy/bilibili-task-docker/exec.log 2>&1 &
```

每天执行一次：10:29

```bash
29 10 * * * /xxxx/yyyy/bilibili-task-docker/exec.sh > /xxxx/yyyy/bilibili-task-docker/exec.log 2>&1 &
```

# 更新镜像

需要手动删除已存在镜像即可

```bash
docker rmi srcrs/bilibili-task:latest
```