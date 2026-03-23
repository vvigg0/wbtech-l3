package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/vvigg0/wbtech-l3/l3/1/internal/handler"
	"github.com/vvigg0/wbtech-l3/l3/1/internal/processor"
	myrabbitmq "github.com/vvigg0/wbtech-l3/l3/1/internal/rabbitmq"
	"github.com/vvigg0/wbtech-l3/l3/1/internal/repository"
	"github.com/vvigg0/wbtech-l3/l3/1/internal/senderTG"
	"github.com/vvigg0/wbtech-l3/l3/1/internal/service"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func Run() error {
	zlog.Init()

	// открываем файл с env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки env файла: %v\n", err)
	}

	//извлекаем все нужные переменные
	dsn := os.Getenv("DB_DSN")
	serverPort := os.Getenv("PORT")
	rabbitURL := os.Getenv("RABBITMQ_URL")
	botToken := os.Getenv("BOT_TOKEN")
	proxy := os.Getenv("PROXY")

	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	db, err := dbpg.New(dsn, nil, opts)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	repo := repository.New(db)

	bot, err := senderTG.New(botToken, proxy)
	if err != nil {
		log.Fatalf("ошибка подключения к боту: %v", err)
	}
	proc := processor.New(bot)
	rabbit := myrabbitmq.New(rabbitURL, proc.HandleMessage)

	service := service.New(repo, rabbit)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		zlog.Logger.Info().Msgf("получен сигнал завершения %v. завершение работы", sig)
		cancel()
	}()

	go func() {
		if err := service.PublishNotifications(ctx); err != nil {
			log.Fatalf("ошибка при публикации уведомлений: %v", err)
		}
	}()
	go func() {
		if err := service.Rabbit.Consumer.Start(ctx); err != nil {
			log.Fatalf("ошибка запуска consumer: %v", err)
		}
	}()

	h := handler.New(service)

	e := ginext.New("")

	registerRoutes(e, h)

	zlog.Logger.Info().Msgf("сервер запущен на %s", serverPort)

	return e.Run(serverPort)
}

func registerRoutes(engine *ginext.Engine, handler *handler.Handler) {
	engine.GET("/notify/:id", handler.GetNotificationStatus)
	engine.POST("/notify", handler.CreateNotification)

	engine.DELETE("/notify/:id", handler.UpdateNotificationStatus)
}
