####Bitcoin-core 节点部署

+ 下载地址
```
wget https://bitcoincore.org/bin/bitcoin-core-0.16.1/bitcoin-0.16.1-x86_64-linux-gnu.tar.gz
```
+ 将包内内容解压⾄至指定⽂文件夹下(绝对路路径,确保路路径存在)
```
tar -zxvf bitcoin-0.16.1-x86_64-linux-gnu.tar.gz -C /hbdata/app
```
+ **节点程序 : /hbdata/app/bitcoin-0.16.1/bin/bitcoind**  
+ **控制台程序 : /hbdata/app/bitcoin-0.16.1/bin/bitcoin-cli**

#### 节点部署 【主网】
>**#rpc账号**  
>rpcuser=bitcoinrpc
>**#rpc密码**  
>rpcpassword=bitcoinrpcpw
>**#rpc连接超时时⻓长**  
>rpctimeout=30
>**#rpc可访问ip**  
>rpcallowip=0.0.0.0/0 #rpc端⼝口(有占⽤用时更更换)
>rpcport=8332 #net连接端⼝口(有占⽤用时更更换)默认不不设置 #port=18332  


### 将配置内容储存为⽂文件(bitcoin.conf)放于节点程序同一目录下  
+ 启动命令:(节点程序⽬录下)  
```
./bitcoind -conf=bitcoin.conf -datadir=/hbdata/app/bitcoin-0.16.1/data -daemon
```
+ 控制台程序启动命令  
```
./bitcoin-cli  -conf=bitcoin.conf -datadir=/hbdata/app/bitcoin-0.16.1/data [opt]
```

