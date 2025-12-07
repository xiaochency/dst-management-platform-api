# 本项目为DMP的旧版本备份
原项目地址：https://github.com/miracleEverywhere/dst-management-platform-api
## :watermelon: 使用方法
>**建议使用 Ubuntu 24系统，低版本系统可能会出现GLIBC版本报错**  
```shell
# 执行以下命令，下载脚本
cd ~ && wget https://dmp-1257278878.cos.ap-chengdu.myqcloud.com/run.sh && chmod +x run.sh
```
```shell
# 自定义启动端口（8082改为你要用的端口），请手动修改run.sh文件或者
sed -i 's/^PORT=.*/PORT=8082/' run.sh
```
```shell
# 根据系统提示输入并回车
./run.sh
```
**安装方法**
```shell
cd ~ && ./run.sh
```
默认启动端口为80，如果您想修改，则修改启动命令：
```shell
# 修改端口为8080
nohup ./dmp -c -l 8080 > dmp.log 2>&1 &
```
## :grapes: 默认用户名密码
>登录后请尽快到右上角用户名-个人中心页面修改密码
>
>>初始密码：
>>admin/123456

---
## :strawberry: 文件介绍
```text
.
├── dmp                 # 主程序
├── dmp.log             # 请求日志
├── dmpProcess.log      # 运行日志
├── DstMP.sdb           # 数据库
├── manual_install.sh   # 饥荒手动安装脚本
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
