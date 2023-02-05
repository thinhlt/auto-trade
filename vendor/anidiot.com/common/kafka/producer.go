package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
)

type Producer struct {
	p sarama.SyncProducer
	topic string
}

func SimpleProducer(brokerList []string, topic string) Producer {
	config := sarama.NewConfig()
	config.Version = sarama.V2_7_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewManualPartitioner
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	return Producer{
		p: producer,
		topic: topic,
	}
}

func (inst Producer) Send(data interface{}) (partition int32, offset int64, err error) {
	dataByte, _ := json.Marshal(data)
	return inst.p.SendMessage(&sarama.ProducerMessage{
		Topic: inst.topic,
		Value: sarama.ByteEncoder(dataByte),
	})
}

func (inst Producer) SendRealSignal(data interface{}) (partition int32, offset int64, err error) {
	dataByte, _ := json.Marshal(data)
	return inst.p.SendMessage(&sarama.ProducerMessage{
		Topic: inst.topic,
		Value: sarama.ByteEncoder(dataByte),
		Partition: 0,
	})
}