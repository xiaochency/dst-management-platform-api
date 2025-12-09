# 本项目为DMP的旧版本备份
原项目地址：https://github.com/miracleEverywhere/dst-management-platform-api
## :watermelon: 使用方法
>**建议使用 Ubuntu 24系统，低版本系统可能会出现GLIBC版本报错**
>**请在root用户下执行本脚本**
```shell
# 执行以下命令，下载脚本
cd ~ && wget https://gh-proxy.com/raw.githubusercontent.com/xiaochency/dst-management-platform-api/refs/heads/main/run.sh && chmod +x run.sh
# 或者
cd ~ && wget  https://ghfast.top/https://raw.githubusercontent.com/xiaochency/dst-management-platform-api/refs/heads/main/run.sh && chmod +x run.sh
```
```shell
# 自定义启动端口（8080改为你要用的端口），请手动修改run.sh文件或者
sed -i 's/^PORT=.*/PORT=8080/' run.sh
```
```shell
# 执行脚本
./run.sh
```
## 安装启动DMP
安装DMP
启动DMP
浏览器访问：http://公网ip:80

## :grapes: 默认用户名密码
>登录后请尽快到右上角用户名-个人中心页面修改密码
>
>>初始密码：
>>admin/123456

## 安装饥荒服务器
选项4：安装饥荒服务器

## 启动服务器
在DMP平台，创建房间或导入存档

---
## :strawberry: 文件介绍
```text
.
├── dmp                 # 主程序
├── dmp.log             # 请求日志
├── dmpProcess.log      # 运行日志
├── DstMP.sdb           # 数据库
└── run.sh              # 运行脚本
```

---

## :peach: 项目介绍
```text
.
├── app
│   ├── auth                    # 登录鉴权
│   ├── externalApi             # 外部接口
│   ├── home                    # 首页
│   ├── logs                    # 日志
│   ├── setting                 # 配置
│   └── tools                   # 工具
├── dist                        # 静态资源
│   ├── assets 
│   ├── index.html
│   ├── index.html.gz
│   └── vite.png
├── docker                      # 容器镜像
│   ├── Dockerfile
│   └── entry-point.sh
├── docs                        # 帮助文档
│   └── images
├── DstMP.sdb                   # 数据库
├── go.mod
├── go.sum
├── LICENSE
├── main.go
├── README.md
├── scheduler                   # 定时任务
│   ├── init.go
│   └── schedulerUtils.go
└── utils                       # 工具集
    ├── constant.go
    ├── exceptions.go
    ├── install.go
    ├── logger.go
    ├── scripts.go
    └── utils.go
```
## 构建
```shell
$env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -o dmp
```

##  :sparkling_heart: 致谢
本项目[前端页面](https://github.com/miracleEverywhere/dst-management-platform-web)基于**koi-ui**二次开发，感谢开源 [@yuxintao6](https://github.com/yuxintao6)  
[[koi-ui gitee]](https://gitee.com/BigCatHome/koi-ui)  
[[koi-ui github]](https://github.com/yuxintao6/koi-ui)  
