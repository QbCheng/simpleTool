package producer

import (
	"github.com/Shopify/sarama"
)

type producerSession struct {
	producer sarama.AsyncProducer
}
