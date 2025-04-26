package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	k "networking/internal/kafka"
)

const (
	topic = "my-topic"
)

// в продакшене, адреса должны браться из переменных окружения
var address = []string{
	"localhost:9092",
}

func main() {
	p, err := k.NewProducer(address)
	if err != nil {
		logrus.Fatal(err)
	}

	// 100 сообщений в топик
	for i := 0; i < 100; i++ {
		textMsg := fmt.Sprintf("Hello World %d", i)
		if err = p.Produce(textMsg, topic); err != nil {
			logrus.Error(err)
		}
	}
}
