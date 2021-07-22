
# Gokins文档

# Gokins: *More Power*

![](https://static01.imgkr.com/temp/5ca8a54f7d6544b6a2c740d5f559e5c4.jpg)




Gokins一款由Go语言和Vue编写的款轻量级、能够持续集成和持续交付的工具.

* **持续集成和持续交付**

  作为一个可扩展的自动化服务器，Gokins 可以用作简单的 CI 服务器，或者变成任何项目的持续交付中心

* **简易安装**

  Gokins 是一个基于 Go 的独立程序，可以立即运行，包含 Windows、Mac OS X 和其他类 Unix 操作系统。


* **安全**

  绝不收集任何用户、服务器信息，是一个独立安全的服务

## Gokins 官网

**地址 : http://gokins.cn**

可在官网上获取最新的Gokins动态

## Gokins Demo
http://gokins.cn:8030
```
用户名: guest
密码: 123456
```

## Quick Start

It is super easy to get started with your first project.


#### Step 1: 环境准备

- Mysql
- Dokcer(非必要)

#### Step 2: 下载
- Linux下载:http://bin.gokins.cn/gokins-linux-amd64
- Mac下载:http://bin.gokins.cn/gokins-darwin-amd64
> 我们推荐使用docker或者直接下载release的方式安装Gokins`

#### Step 3: 启动服务

```
./gokins
``` 

#### Step 3: 安装Gokins

访问 `http://localhost:8030`进入到Gokins安装页面

![](https://static01.imgkr.com/temp/e484d9747dec43108325c22283abe39f.png)

按页面上的提示填入信息

默认管理员账号密码

`username :gokins `

`pwd: 123456 `

#### Step 4:  新建流水线

- 进入到流水线页面

![](https://static01.imgkr.com/temp/ce383350056d4a63872b868c8f169c39.png)



- 点击新建流水线

![](https://static01.imgkr.com/temp/a3c2a870c9d94956bda2a685cc447077.png)


填入流水线基本信息

- 流水线配置

```
version: 1.0
vars:
stages:
  - stage:
    displayName: build
    name: build
    steps:
      - step: shell@sh
        displayName: test-build
        name: build
        env:
        commands:
          - echo Hello World

```

关于流水线配置的YML更多信息请访问 [YML文档](http://gokins.cn/%E5%B7%A5%E4%BD%9C%E6%B5%81%E8%AF%AD%E6%B3%95/)


- 运行流水线

![](https://static01.imgkr.com/temp/f002a22738644c8dbd40f0860c2bbb9e.png)


`这里可以选择输入仓库分支或者commitSha,如果不填则为默认分支`

- 查看运行结果

![](https://static01.imgkr.com/temp/681c8ea0a7dc45bcb9fe14234c5761be.png)

