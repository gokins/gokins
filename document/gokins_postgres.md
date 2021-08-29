# gokins+postgres安装示例说明

环境整体依赖于docker进行配置说明

#### Step 1: 环境准备

```shell script
# 启动postgres
sudo docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres

# 然后在postgres中新建数据库`gokins`

```

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

数据库地址填写postgres的连接地址：`127.0.0.1:5432`，数据库名称：`gokins`，数据库用户：`postgres`，数据库密码：`password`。***实际配置信息请根据自己的实际环境信息进行填写***

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
