messaging
=========

IM项目框架

## 基础环境

- 包管理工具 [dep](https://github.com/golang/dep)
- 协议格式化工具 [protoc](https://github.com/google/protobuf)
- 协议插件 [protoc-gen-go](https://github.com/golang/protobuf)
- 微服务协议插件 [protoc-gen-micro](https://github.com/micro/protoc-gen-micro)

## 项目文件

- conf 配置文件
- script 服务创建，编译的脚本
- src/lib 基础库
- src/gateway tcp网关，提供长连接服务，默认20000端口

## 运行

### 下载依赖包

src目录下，运行：

```
dep ensure -v
```

### 编译

script目录下，运行：
```
bash build.sh
```
生成msgbin目录，编译完成的服务会放在此目录下

## 创建服务

script目录下，运行：
```
bash createservice.sh xxx
```
脚本利用template服务模板创建新的服务，xxx为需要创建的服务名称

### 模板服务说明

- client 其他服务调用该服务的客户端
- proto 服务协议文件
- server 服务端
- test 测试用例
- main.go 服务入口