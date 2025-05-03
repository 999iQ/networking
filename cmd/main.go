package main

import (
	"gorm.io/gorm"
	"networking/internal/config"
	"networking/internal/database"
	"networking/internal/logger"
)

//const (
//	topic = "my-topic"
//)

var db *gorm.DB

func main() {

	appLog := logger.New()

	cfg, err := config.LoadConfig()
	if err != nil {
		appLog.Fatalf("ошибка при загрузке конфигурации: %v", err)
	}

	db, err = database.ConnectDB(cfg, appLog)
	if err != nil {
		appLog.Fatalf("ошибка при подключении к базе данных: %v", err)
	}
	defer database.CloseDB(db)

	appLog.Info("Приложение запущено!")

}

//func main() {
//p, err := kafka.NewProducer(address)
//if err != nil {
//	logrus.Fatal(err)
//}
//
//// 100 сообщений в топик
//for i := 0; i < 100; i++ {
//	textMsg := fmt.Sprintf("Hello World %d", i)
//	if err = p.Produce(textMsg, topic); err != nil {
//		logrus.Error(err)
//	}
//}

// test db

//}
