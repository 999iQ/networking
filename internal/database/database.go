package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"networking/internal/config"
	"networking/internal/logger"
)

func ConnectDB(cfg config.DbConnectConfig, log *logger.Logger) (*gorm.DB, error) {
	connectStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname, cfg.Sslmode)

	db, err := gorm.Open(postgres.Open(connectStr), &gorm.Config{
		Logger: log.GormLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	log.Info("Успешное подключение к базе данных!")

	// migrate check
	if cfg.AutoMigrate {
		log.Info("Выполняется автоматическая миграция...")
		// my new db models
		//err = db.AutoMigrate(&models...)
		//if err != nil {
		//	return nil, fmt.Errorf("ошибка при автоматической миграции: %w", err)
		//}
		log.Info("Автоматическая миграция выполнена.")
	} else {
		log.Info("Автоматическая миграция отключена (AUTO_MIGRATE != true).")
	}
	return db, nil

}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("ошибка при получении *sql.DB: %v", err)
	}
	if err = sqlDB.Close(); err != nil {
		log.Fatalf("ошибка при закрытии соединения с БД: %v", err)
	}
}
