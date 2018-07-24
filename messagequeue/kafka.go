package messagequeue

import (
	"github.com/Shopify/sarama"
	"github.com/gembackend/gembackendlog"
	"github.com/gembackend/conf"
)

var (
	log = gembackendlog.Log

	msg = &sarama.ProducerMessage{
		Partition: int32(-1),
		Key:       sarama.StringEncoder("key"),
	}

	kafkaurl string
)

func init() {
	kafkaurl = conf.MysqlHost + ":" + conf.KafkaPort
}

func MakeProducer() (producer sarama.SyncProducer) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	producer, err := sarama.NewSyncProducer([]string{kafkaurl}, config)

	if err != nil {
		panic(err)
	}
	return
}

func MakeMessage(topic, value string, producer sarama.SyncProducer) {
	msg.Value = sarama.ByteEncoder(value)
	msg.Topic = topic
	paritition, offset, err := producer.SendMessage(msg)

	if err != nil {
		log.Errorf("Send Message Fail %s", err)
	}
	log.Infof("Partion = %d, offset = %d\n", paritition, offset)
}

func MakeConsumer() (consumer sarama.Consumer) {
	consumer, err := sarama.NewConsumer([]string{kafkaurl}, nil)

	if err != nil {
		panic(err)
	}

	return
}

func MakePcs(consumer sarama.Consumer, topic string) (pcs []sarama.PartitionConsumer) {
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		panic(err)
	}
	pcs = make([]sarama.PartitionConsumer, 0)
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}
		pcs = append(pcs, pc)
	}
	return
}
