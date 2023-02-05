package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
	"sync"
)

func SimpleClient(brokerList []string, groupId string) (client sarama.ConsumerGroup, err error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_7_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	client, err = sarama.NewConsumerGroup(brokerList, groupId, config)
	return
}

func StartListen(client sarama.ConsumerGroup, topics []string, consumer Consumer, stopChan chan bool)  {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, topics, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-stopChan:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err := client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}


// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	Ready   chan bool
	Handler func([]byte)
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as Ready
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s, partition = %d", string(message.Value), message.Timestamp, message.Topic, message.Partition)
		consumer.Handler(message.Value)
		session.MarkMessage(message, "")
	}

	return nil
}
