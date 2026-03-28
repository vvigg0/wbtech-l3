package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	dsn := fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"))

	botToken := os.Getenv("BOT_TOKEN")
	proxy := os.Getenv("PROXY")

	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	db, err := dbpg.New(dsn, nil, opts)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}
	defer db.Master.Close()

	repo := repository.New(db)

	bot, err := senderTG.New(botToken, proxy)
	if err != nil {
		return fmt.Errorf("ошибка подключения к боту: %w", err)
	}

	proc := processor.New(bot, repo)

	rabbitURL := "amqp://guest:guest@rabbit:5672/"
	rabbit, err := myrabbitmq.New(rabbitURL, proc.HandleMessage)
	if err != nil {
		return fmt.Errorf("ошибка инициализации rabbit: %w", err)
	}
	defer rabbit.Client.Close()

	service := service.New(repo, rabbit.Publisher)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigs)

	errorCh := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Go(func() {
		if err := service.PublishNotifications(ctx); err != nil && !errors.Is(err, ctx.Err()) {
			select {
			case errorCh <- fmt.Errorf("ошибка при публикации уведомлений: %w", err):
			default:
			}
		}
	})
	wg.Go(func() {
		if err := rabbit.Consumer.Start(ctx); err != nil && !errors.Is(err, ctx.Err()) {
			select {
			case errorCh <- fmt.Errorf("ошибка запуска consumer: %w", err):
			default:
			}
		}
	})

	h := handler.New(service)

	router := ginext.New("")
	registerRoutes(router, h)
	router.Static("/static", "./web")
	router.LoadHTMLFiles("./web/index.html")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router}

	wg.Go(func() {
		zlog.Logger.Info().Msgf("сервер запущен")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			select {
			case errorCh <- fmt.Errorf("ошибка сервера: %w", err):
			default:
			}
		}
	})

	var runErr error
	select {
	case runErr = <-errorCh:
	case sig := <-sigs:
		zlog.Logger.Info().Msgf("получен сигнал завершения %v. завершение работы", sig)
	}
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		if runErr == nil {
			runErr = fmt.Errorf("принужденное завершение сервера: %w", err)
		}
	}

	wg.Wait()
	return runErr
}

func registerRoutes(engine *ginext.Engine, handler *handler.Handler) {
	engine.GET("/", func(ctx *ginext.Context) {
		ctx.HTML(200, "index.html", nil)
	})
	engine.GET("/notify/:id", handler.GetNotificationStatus)
	engine.POST("/notify", handler.CreateNotification)

	engine.DELETE("/notify/:id", handler.CancelNotification)
}
