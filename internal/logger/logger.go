package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Logger struct {
	*logrus.Logger
	GormLogger logger.Interface
}

func New() *Logger {
	log := logrus.New()
	// settings format log
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true, // Чтобы timestamp отображался полностью
	})

	// Устанавливаем уровень логирования из переменной окружения LOG_LEVEL
	logLevelStr := os.Getenv("LOG_LEVEL")
	logLevel, err := logrus.ParseLevel(logLevelStr)
	if err != nil {
		// Если LOG_LEVEL не установлен или некорректен, используем INFO по умолчанию
		logLevel = logrus.InfoLevel
		log.Printf("LOG_LEVEL не установлен или некорректен, используем INFO по умолчанию") // Используем Logrus Printf для начальной загрузки
	}
	log.SetLevel(logLevel)

	// Настройка вывода (stdout и/или файл)
	output := os.Getenv("LOG_OUTPUT")
	if output == "file" {
		// Убедитесь, что директория для логов существует
		logDir := os.Getenv("LOG_DIR")
		if logDir == "" {
			logDir = "." // Текущая директория по умолчанию
		}
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			err := os.MkdirAll(logDir, 0755) // Создаем директорию, если ее нет
			if err != nil {
				log.Errorf("Не удалось создать директорию для логов: %v", err)
				log.SetOutput(os.Stdout) // Переключаемся на вывод в stdout
			}
		}

		logFile, err := os.OpenFile(logDir+"/app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Errorf("Не удалось открыть файл логов: %v", err)
			log.SetOutput(os.Stdout) // Переключаемся на вывод в stdout
		} else {
			log.SetOutput(io.MultiWriter(os.Stdout, logFile)) // output in file & terminal
		}
	} else {
		log.SetOutput(os.Stdout) // output in terminal
	}

	// Создание GORM логгера
	gormLogger := NewGormLogger(log)

	return &Logger{
		Logger:     log,
		GormLogger: gormLogger,
	}
}

// NewGormLogger создает GORM логгер, использующий Logrus.
func NewGormLogger(log *logrus.Logger) logger.Interface {
	return &gormLogger{
		Log: log,
	}
}

type gormLogger struct {
	Log *logrus.Logger
}

// LogMode реализация интерфейса gorm.logger.Interface
func (g *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *g
	return &newlogger
}

// Info реализация интерфейса gorm.logger.Interface
func (g *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	g.Log.WithContext(ctx).Infof(msg, data...)
}

// Warn реализация интерфейса gorm.logger.Interface
func (g *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	g.Log.WithContext(ctx).Warnf(msg, data...)
}

// Error реализация интерфейса gorm.logger.Interface
func (g *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	g.Log.WithContext(ctx).Errorf(msg, data...)
}

// Trace реализация интерфейса gorm.logger.Interface
func (g *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && g.Log.GetLevel() >= logrus.ErrorLevel:
		sql, rows := fc()
		g.Log.WithContext(ctx).Errorf("%s [%s] [%v] [%v] %s", utils.FileWithLineNum(), err, elapsed, rows, sql)
	case elapsed > time.Second && g.Log.GetLevel() >= logrus.WarnLevel:
		sql, rows := fc()
		g.Log.WithContext(ctx).Warnf("%s [%s] [%v] [%v] %s", utils.FileWithLineNum(), elapsed, rows, sql)
	case g.Log.GetLevel() == logrus.DebugLevel:
		sql, rows := fc()
		g.Log.WithContext(ctx).Debugf("%s [%s] [%v] [%v] %s", utils.FileWithLineNum(), elapsed, rows, sql)
	}
}
