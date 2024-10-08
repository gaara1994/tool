/*
Copyright © 2024 yantao
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
var kafkatestCmd = &cobra.Command{
	Use:     "kafkaTest",
	Aliases: []string{"kt"},
	Short:   "测试kafka (别名: kt)",
	Long: `
连接测试kafka。

使用方法:
  tool kafkaTest [broker] [topic]
  tool kt 		 [broker] [topic]  # 使用别名
`,
	Run: kafkaTest,
}

func init() {
	rootCmd.AddCommand(kafkatestCmd)
}
func kafkaTest(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: tool tk broker topic")

		return
	}

	addr := args[0] // 从命令行参数中获取端口号
	fmt.Println("kafka 地址为：", addr)

	topic := "test" // 默认主题
	if len(args) > 1 {
		topic = args[1] // 如果提供了主题，则使用提供的主题
	}
	fmt.Println("Kafka 主题为：", topic)

	// 确保主题存在
	if err := ensureTopicExists(addr, topic); err != nil {
		log.Fatalf("Failed to ensure topic exists: %v", err)
	}

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
			fmt.Printf("Consumer error: %v\n", err)
		case <-ctx.Done():
			fmt.Println("Context timeout, stopping consumer")
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

func ensureTopicExists(addr, topic string) error {
	// 创建一个新的 Sarama 配置
	config := sarama.NewConfig()
	config.Version = sarama.V3_0_0_0 // 根据你的 Kafka 版本设置

	// 创建一个 ClusterAdmin 实例
	admin, err := sarama.NewClusterAdmin([]string{addr}, config)
	if err != nil {
		return fmt.Errorf("failed to create cluster admin: %v", err)
	}
	defer admin.Close()

	// 尝试描述主题以检查其是否存在
	_, err = admin.DescribeTopics([]string{topic})
	if err != nil {
		// 如果找不到主题，则尝试创建它
		if kerr, ok := err.(sarama.KError); ok && kerr == sarama.ErrUnknownTopicOrPartition {
			// 定义新主题的详细信息
			topicDetail := &sarama.TopicDetail{
				NumPartitions:     1, // 设置分区数量
				ReplicationFactor: 1, // 设置副本因子
			}

			// 创建新主题
			err = admin.CreateTopic(topic, topicDetail, false) // 最后一个参数是是否等待所有副本都可用
			if err != nil {
				return fmt.Errorf("failed to create topic: %v", err)
			}
			fmt.Printf("Topic %s created successfully.", topic)
		} else {
			return fmt.Errorf("failed to describe topics: %v", err)
		}
	} else {
		fmt.Printf("Topic %s already exists.\n", topic)
	}

	return nil
}
