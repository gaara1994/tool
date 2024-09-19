/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/spf13/cobra"
)

// testKafkaCmd represents the testKafka command
var testKafkaCmd = &cobra.Command{
	Use:     "kafka",
	Aliases: []string{"tk", "kafka"},
	Short:   "测试kafka",
	Long:    ``,
	Run:     testKafka,
}

func init() {
	rootCmd.AddCommand(testKafkaCmd)
}
func testKafka(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: tool tk addr topic")
		return
	}

	addr := args[0] // 从命令行参数中获取端口号
	fmt.Println("kafka 地址为：", addr)

	topic := "test" // 默认主题
	if len(args) > 1 {
		topic = args[1] // 如果提供了主题，则使用提供的主题
	}
	fmt.Println("Kafka 主题为：", topic)

	// 创建一个等待组，用于等待所有协程完成
	var wg sync.WaitGroup

	// 开启消费者协程
	wg.Add(1)
	go func() {
		defer wg.Done()
		startConsumer(addr, topic)
	}()

	// 开启生产者协程
	wg.Add(1)
	go func() {
		defer wg.Done()
		startProducer(addr, topic)
	}()

	// 等待所有协程完成
	wg.Wait()

}

func startConsumer(addr, topic string) {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true
	consumerConfig.Version = sarama.V3_0_0_0 // 根据你的 Kafka 版本设置

	consumer, err := sarama.NewConsumer([]string{addr}, consumerConfig)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest) // 从最新的偏移量开始消费
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Consumer: Received message: %s\n", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			log.Printf("Consumer error: %v\n", err)
		case <-ctx.Done():
			log.Println("Context timeout, stopping consumer")
			return
		}
	}
}

func startProducer(addr, topic string) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Version = sarama.V3_0_0_0 // 根据你的 Kafka 版本设置

	producer, err := sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder("Hello, Kafka!"),
	}

	for {
		partition, offset, err := producer.SendMessage(message)
		if err != nil {
			log.Fatalf("Failed to send message: %v", err)
		}
		fmt.Printf("Producer: Message sent to topic %s, partition %d at offset %d\n", topic, partition, offset)

		// 为了演示，发送一条消息后退出
		time.Sleep(1 * time.Second) // 等待一段时间以确保消息被消费者读取
	}

}
