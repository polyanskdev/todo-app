package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/polyanskdev/todo-app"
	"github.com/polyanskdev/todo-app/pkg/handler"
	"github.com/polyanskdev/todo-app/pkg/repository"
	"github.com/polyanskdev/todo-app/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run() {
	initLogger()
	initConfig()
	initEnv()
	initApplication()
}

func initLogger() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
}

func initConfig() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Ошибка чтения файлов конфигурации: %s", err.Error())
	}
}

func initEnv() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Ошибка инициализации переменных окружения: %s", err.Error())
	}
}

func initApplication() {
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("Ошибка инициализации БД: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Ошибка запуска http сервера: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Ошибка завершения работы сервера: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("Ошибка при закрытии соединения с БД: %s", err.Error())
	}
}
