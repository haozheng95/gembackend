####kafka
+ kafka docker镜像
```cgo
docker pull spotify/kafka    
docker run -p 2181:2181 -p 9092:9092 --env ADVERTISED_HOST=39.106.55.83 --env ADVERTISED_PORT=9092 -d --name kafka spotify/kafka
export KAFKA=39.106.55.83:9092
/opt/kafka_2.11-0.10.1.0/bin/kafka-console-producer.sh --broker-list $KAFKA --topic test
export ZOOKEEPER=39.106.55.83:2181
/opt/kafka_2.11-0.10.1.0/bin/kafka-console-consumer.sh --zookeeper $ZOOKEEPER --topic test
```
+ 创建Topic
```cgo
bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test
```
+ 列出Topic是否创建成功
```cgo
bin/kafka-topics.sh --list --zookeeper localhost:2181
```
+ 发送消息 向创建的test Topic 发送消息(生产者)
```cgo
bin/kafka-console-producer.sh --broker-list localhost:9092 --topic test
```
+ 创建消费者 订阅一个test Topic,并进行消费  
```cgo
bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test --from-beginning
```

PS:如果你的代码不是运行在loalhost这台机器上的话,需要修改 config/server.properties 配置文件的listeners中的host,否则kafka服务端会拒绝你非localhost的连接请求,配置好后重启kafka服务.
