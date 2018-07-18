#### Parity eth节点安装

* 节点程序下载  
```parity_1.9.5_ubuntu_amd64.deb```

* wget 下载地址  
``` 
wget http://d1h4xl4cr1h0mo.cloudfront.net/v1.9.5/x86_64-unknown-linux-gnu/parity_1.9.5_ubuntu_amd64.deb
```

* 将包内内容解压⾄至指定⽂文件夹下(确保路径存在)
```
dpkg -X ./parity_1.9.5_ubuntu_amd64.deb /your/path/parity
```

#### 控制台(使⽤用geth节点程序登录parity控制台)

* 下载release版的geth程序
```
wget https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.8.11-dea1ce05.tar.gz
```

* 解压到指定目录

```commandline. 
tar -zxvf geth-linux-amd64-1.8.3-329ac18e.tar.gz -C /your/path/geth
```

#### 节点主网部署

* 如果不指定节点配置文件路径，节点默认读取
~/.local/share/io.parity.ethereum/目录下的配置

* 写一个配置文件config.toml来指定节点基本配置
```commandline
vim /your/path/config.toml
```

* 基本配置
>[parity]  
>base_path = "/hbdata/app/parity/usr/bin/data"  
>db_path = "/hbdata/app/parity/usr/bin/data/chains"  
>keys_path = "/hbdata/app/parity/usr/bin/data/keys"
* 网络配置
>[network]
port = 30300

* rpc配置
>[rpc]  
>prot = 8542  
>interface = "all"  
>apis = ["web3", "eth", "net", "parity", "traces", "parity_set", "rpc", "personal"]


#### 启动主网
```commandline
nohup ./parity --config config.toml >parity.log 2>&1 &
```

#### 输出日志监控
```commandline
tail -f parity.log
```

#### 添加节点

* 把节点地址分行存到nodes.log文件中 再重新启动
```commandline
nohup ./parity --config config.toml --reserved-peers nodes.log >parity.log 2>&1 &
```

#### geth链接主网
```commandline
geth attach http://127.0.0.1:8542
```
