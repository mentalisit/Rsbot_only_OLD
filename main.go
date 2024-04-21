package main

import (
	"Rsbot_only/internal/bot"
	"Rsbot_only/internal/clients"
	"Rsbot_only/internal/config"
	"Rsbot_only/internal/storage"
	"fmt"
	"github.com/mentalisit/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Bot loading ")

	err := RunNew()
	if err != nil {
		fmt.Println("Error loading bot", err)
		time.Sleep(10 * time.Second)
		panic(err.Error())
	}
}

func RunNew() error {
	//читаем конфигурацию с ENV
	cfg := config.InitConfig()

	//создаем логгер
	log := logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook)

	if cfg.BotMode == "dev" {
		log = logger.LoggerZapDEV()

		//os.Exit(1)
		go func() {
			time.Sleep(5 * time.Minute)
			os.Exit(1)
		}()
	}

	log.Info("🚀  загрузка  🚀 " + cfg.BotMode)

	//storage
	st := storage.NewStorage(log, cfg)

	//clients Discord, Telegram
	cl := clients.NewClients(log, st, cfg)
	go bot.NewBot(st, cl, log, cfg)

	//ожидаем сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	return nil
}
