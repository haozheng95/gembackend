# gembackend-go
#### go版本钱包 支持erc20的token转账 btc1地址转账  
+ 项目分成server部分和script部分
+ server部分负责提供api调用 发送交易 查看余额 查看交易记录等。
+ script部分负责从节点将相关地址的数据从节点上同步下来
oo
## 添加docker file
```cgo
docker pull yinhaozheng/gembackend-go
```

### 脚本启动
```cgo
cd run
go build
./run -action=eth-kafka-script #启动kafka队列脚本 
./run -action=eth-updater-web3-mul #启动多线程web3更新程序
./run -action=eth-updater-web3 #启动单线程web3更新程
./run -action=eth-updater-ethscan #启动ethscan接口的更新程序
```

### 配置文件目录在conf下
+ 启动需要修改你自己的节点，数据库，和kafka配置
