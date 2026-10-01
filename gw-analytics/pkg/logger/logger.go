package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func LoggerInit(service string) {
	// дописать что бы было Id транзакции

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	log, _ = cfg.Build(
		zap.Fields(
			zap.String("Сервис", service)),
	)

}

func Debug(msg string, fields ...zap.Field) { log.Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { log.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { log.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { log.Error(msg, fields...) }

func Sync() { _ = log.Sync() }

/*
Не логируй там, где прокидываешь ошибку наверх (внутри handle можно просто return fmt.Errorf(...),
а логировать в Run).

Логируй на границе — где ошибка превращается в решение «continue / return / ретрай».

Уровни: context.Canceled при shutdown — это Info (штатное завершение).
Ошибка Kafka — Error.Ошибка парсинга одного сообщения — Error + continue.

*/
