package kafka

import (
	"game-server/pkg/logger"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

var Producer sarama.SyncProducer

func Init() {
	// 暂时硬编码Kafka地址，绕过配置读取问题
	kafkaAddr := "localhost:9092"
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	producer, err := sarama.NewSyncProducer([]string{kafkaAddr}, config)
	if err != nil {
		logger.Error("Kafka连接失败", zap.Error(err))
		// 暂时不panic，允许服务启动
		logger.Warn("Kafka连接失败，服务将继续运行，但匹配功能可能受限")
		return
	}
	Producer = producer
	logger.Info("Kafka初始化成功")
}

// 发送消息
func Send(topic string, msg string) {
	if Producer == nil {
		logger.Warn("Kafka未初始化，消息发送失败", zap.String("topic", topic), zap.String("msg", msg))
		return
	}

	partition, offset, err := Producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	})
	if err != nil {
		logger.Error("Kafka发送失败", zap.Error(err), zap.String("topic", topic))
	} else {
		logger.Info("Kafka发送成功", zap.Int32("partition", partition), zap.Int64("offset", offset), zap.String("topic", topic))
	}
}
