## chatgpt机器人（qq版）

### 工具

- python3.8
- golang 1.19 + echo + requests
- mysql 8.20
- **docker 20.2**
- 阿里云香港服务器 2c2g

- **chatgpt** **api**

- **cqhttp**（基于miraiGo）



### 功能

- [x] chatgpt 接入qq对话
- [x] 官方的记忆功能（官方api限制总字数）
- [x] 清空数据功能
- [x] 授予某用户权限
- [x] 私聊功能（需要管理员添加权限）
- [x] 架构完善和改进
- [x] 预设人设功能，可以添加人设（不完善）
- [ ] 多线程，集成

## 使用须知
1. 在model下创建**conf.ymal**，用来连接数据库
```
    sql:
    username: [用户名]
    password: [密码]
    db_name: [数据库名称]
```

2. 在utils下创建**basicData.go** ,用于基本数据
```
var qqServe QqServe = Deal{}
var qqUid = "你的qq号"         //填
var AtQqUid = fmt.Sprintf("[CQ:at,qq=%v]", qqUid)
var LenAtQqUid = len(AtQqUid) + 1
var ip = "127.0.0.1"   //填你的部署的ip地址，这里是本地
var setting = model.AllSetting
``` 

## 一图看懂架构

![](https://raw.githubusercontent.com/Solituderr/goChatQqBot/master/images/1241.gif)




### 说明

目前作者还未集成所有代码，需要使用者自己部署到服务器（本项目纯docker部署）。

此项目仅供交流使用，不提供相关操作流程，不提供release，需要自己部署。

所有操作与作者无关，造成的相关影响由使用者自负。

