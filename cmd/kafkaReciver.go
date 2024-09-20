/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/spf13/cobra"
)

// kafkaReciverCmd represents the kafkaReciver command
var kafkaReciverCmd = &cobra.Command{
	Use:     "kafkaReciver",
	Aliases: []string{"kr"},
	Short:   "接收kafka消息 (别名: kr)",
	Long: `
接收kafka消息。

使用方法:
  tool kafkaReciver [broker] [topic] 
  tool kr 		    [broker] [topic]  # 使用别名
  `,
	Run: kafkaReciver,
}

func init() {
	rootCmd.AddCommand(kafkaReciverCmd)
}

func kafkaReciver(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: tool kr broker topic")
		return
	}
	broker := args[0]
	topic := args[1]
	fmt.Println("broker:", broker)
	fmt.Println("topic:", topic)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V3_0_0_0

	consumer, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		fmt.Println("Failed to start consumer: ", err)
		return
	}
	defer consumer.Close()

	// 获取partition列表
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Printf("Failed to get partitions for topic %s: %v", topic, err)
		return
	}

	// 对于每个partition创建一个goroutine来消费消息
	for _, partition := range partitions {
		go func(partition int32) {
			pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				log.Fatalf("Failed to start consumer for partition %d: %v", partition, err)
			}
			defer pc.AsyncClose()

			// 消费消息
			for msg := range pc.Messages() {
				fmt.Printf("Partition: %d, Offset: %d, Key: %s, Value: %s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(partition)
	}

	// 处理SIGINT和SIGTERM信号以优雅地关闭消费者
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	fmt.Println("Shutting down consumer...")
}
