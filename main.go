package main

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TextFilter(msg *gotgbot.Message) bool {
	return msg.Text != ""
}

var logger *zap.Logger

func main() {
	// Configuring zap logger
	config := zap.NewProductionConfig()
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var err error
	logger, err = config.Build()
	if err != nil {
		panic("Failed to initialize logger")
	}

	// Creating a new bot instance
	bot, err := gotgbot.NewBot("YOUR_BOT_API_KEY", nil)
	if err != nil {
		logger.Sugar().Panicf("Failed to create bot: %v", err)
	}
	logger.Sugar().Info("Bot Started")

	// Create a dispatcher with default options
	dispatcher := ext.NewDispatcher(nil)

	// Add echo handler
	dispatcher.AddHandler(handlers.NewMessage(TextFilter, echo))

	// Create an updater and start polling
	updater := ext.NewUpdater(dispatcher, nil)
	err = updater.StartPolling(bot, &ext.PollingOpts{})
	if err != nil {
		logger.Sugar().Panicf("Failed to start polling: %v", err)
	}

	updater.Idle()
}

// echo is a handler function that echoes back the received message
func echo(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := b.SendMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.Text, nil)
	if err != nil {
		logger.Sugar().Errorf("Failed to send message: %v", err)
	}
	return err
}
