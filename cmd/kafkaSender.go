/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/spf13/cobra"
)

// kafkaSenderCmd represents the kafkaSender command
var kafkaSenderCmd = &cobra.Command{
	Use:     "kafkaSender",
	Short:   "发送消息到kafka (别名: ks)",
	Aliases: []string{"ks"},
	Long: `
发送消息到kafka。

使用方法:
  tool kafkaSender [broker] [topic] [jsonMessage]
  tool ks 		   [broker] [topic] [jsonMessage] # 使用别名
	`,
	Run: kafkaSender,
}

func init() {
	rootCmd.AddCommand(kafkaSenderCmd)
}

func kafkaSender(cmd *cobra.Command, args []string) {
	if len(args) < 3 {
		fmt.Println("参数不足")
		fmt.Println("Usage: tool ks [broker] [topic] [jsonMessage]")
		return
	}

	var broker = args[0]
	var topic = args[1]
	var jsonMessage = args[2]

	fmt.Println("broker:", args[0])
	fmt.Println("topic:", args[1])
	fmt.Println("jsonMessage:", args[2])

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Version = sarama.V3_0_0_0
	producer, err := sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		fmt.Println("Failed to create producer", err)
		return
	}
	defer producer.Close()

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonMessage),
	}
	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		fmt.Println("Failed to send message", err)
		return
	}
	fmt.Printf("Producer: Message sent to topic %s, partition %d at offset %d\n", args[1], partition, offset)
	fmt.Println("发送完成")
}
