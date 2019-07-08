## 简介
一个基于Gin框架的封装轻量级Go语言业务框架，代码编写上分成Controller Model Service 三层（目前不支持返回视图）。
框架目录结构参考PHP Laravel框架，易于让习惯了PHP MVC WEB框架开发的同学上手。

## 开发环境搭建

### 修改环境变量

```
sudo vim ~/.bash_profile
export GOPATH="/go"
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
export GO111MODULE=on
export GOPROXY=https://goproxy.io
source ~/.bash_profile
```

### 加载依赖
```
go mod vendor
```

### Docker
```text
1. echo "127.0.0.1 api-demo.qd-docker.com >> /etc/hosts"
2. cp .env.example .env # 配置文件
3. cd build/docker
4. docker-compose up --force-recreate -d # 启动容器
5. curl "http://api-demo.qd-docker.com/hello" # 访问测试
6. docker exec -it go.demo sh # 进入容器
```

### 本机
```text
1. cp .env.example .env # 配置文件
2. sh build/shell/build.sh  #编译
3. build/bin/snow -a api #启动Api服务
4. build/bin/snow -a cron #启动Cron定时任务服务
5. build/bin/snow -a job #启动Job队列服务  参数 -queue "topic1,topic2" 或者 --queue="topic1,topic2"
4. curl "http://127.0.0.1:8000/hello" #访问测试
```
    
## 新项目
使用build/shell/replace.sh批量本地项目文件的包命名空间
```sh
sh build/shell/replace.sh [新项目的包命名空间，如github.com/qit-team/snow]
```

## Wiki
[https://github.com/qit-team/snow/wiki](https://github.com/qit-team/snow/wiki)



