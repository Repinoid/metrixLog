package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

var (
	Logger  *slog.Logger
	pollInt int = 1
)

func main() {
	
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println("Couldn't load .env ", err)
		// It's important to note that it WILL NOT OVERRIDE an env variable
		// that already exists - consider the .env file to set dev vars or sensible defaults.
		// Не прерываем выполнение, так как переменные могут быть установлены в окружении
	}
	
	// определяем уровень логирования из переменной окружения, заданной в docker-compose.yml или в .env
	level := slog.LevelInfo
	switch GetEnv("LOG_LEVEL", "INFO") {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: false, //NO  Добавлять информацию об исходном коде
	})
	Logger = slog.New(handler)
	slog.SetDefault(Logger)

	// POLL_DURATION скважность логирования метрик
	pollInt, err = strconv.Atoi(GetEnv("POLL_DURATION", "1"))
	if err != nil {
		Logger.Error("Bad strconv.Atoi (POLL_DURATION) ", "", err)
		return
	}

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
		<-exit
		cancel()
	}()

	var wg sync.WaitGroup

	const chanCap = 4

	metroBarn := make(chan map[string]uint64, chanCap)

	wg.Add(2)
	go metrixIN(ctx, metroBarn, &wg)
	go bolda(ctx, metroBarn, &wg)

	Logger.Info("Goroutines started")

	wg.Wait()
	close(metroBarn)
	Logger.Info("Agent Shutdown gracefully")
	return nil
}

// получает банчи метрик и складывает в barn
// func metrixIN(ctx context.Context, metroBarn chan<- []models.Metrics, wg *sync.WaitGroup, sigint chan os.Signal) {
func metrixIN(ctx context.Context, metroBarn chan<- map[string]uint64, wg *sync.WaitGroup) {
	defer wg.Done()
	// var memStorage map[string]uint64
	// var err error
	tickerPoll := time.NewTicker(time.Duration(pollInt) * time.Second)
	for {
		select {
		case <-ctx.Done():
			Logger.Info("Горутина metrixIN остановлена")
			return
		// по тикеру запрашиваем метрикис рантайма
		case <-tickerPoll.C:
			memStorage, err := GetMetrixFromOS()
			if err != nil {
				return
			}
			// засылаем метрики в канал
			metroBarn <- memStorage
		}
	}
}

// работник отсылает банчи метрик на сервер,
func bolda(ctx context.Context, metroBarn <-chan map[string]uint64, wg *sync.WaitGroup) {
	defer wg.Done()
	var bunch map[string]uint64
	for {
		select {
		case <-ctx.Done():
			Logger.Info("Горутина bolda остановлена")
			return
		case bunch = <-metroBarn:
			sendMetrics(bunch)
		}
	}
}

func sendMetrics(bunch map[string]uint64) {

	Logger.Info("SysMetrix logs",
		slog.Any("metrics", bunch))

}

// GetEnv возвращает значение переменной окружения или значение по умолчанию
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
