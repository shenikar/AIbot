package logger

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger(logFile string) (*Logger, error) {
	log := logrus.New()

	// Форматирование
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		ForceColors:     true, // Цветные логи в консоли
	})

	// Открываем файл для записи логов
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// Вывод логов в консоль и в файл одновременно
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)

	log.SetLevel(logrus.InfoLevel) // Уровень логирования

	return &Logger{log}, nil
}
