package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/configs"
	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/pressly/goose/v3"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configFile string

var migrationsPath string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
	flag.StringVar(&migrationsPath, "migrations", "migrations", "Path to migrations directory")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	configData, err := config.NewConfigFromFile(configFile)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	logg := logger.New(configData.Logger.Level)

	var storage app.Storage
	switch configData.Storage.Type {
	case "memory":
		storage = memorystorage.New()
	case "sql":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			configData.DB.Host, configData.DB.Port, configData.DB.User, configData.DB.Password, configData.DB.DBName)

		if err := runMigrations(dsn); err != nil {
			panic("failed to apply migrations: " + err.Error())
		}

		storage, err = sqlstorage.New(dsn)
		if err != nil {
			panic("failed to connect to db: " + err.Error())
		}
	default:
		panic("unknown storage type: " + configData.Storage.Type)
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, configData.Server.Host, configData.Server.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func runMigrations(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	fmt.Printf("DSN: %s\n", dsn)
	fmt.Printf("Migrations path: %s\n", migrationsPath)

	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		fmt.Printf("Error reading migrations directory: %v\n", err)
	} else {
		fmt.Printf("Files in migrations directory:\n")
		for _, file := range files {
			fmt.Printf("  - %s\n", file.Name())
		}
	}

	if err := goose.Up(db, migrationsPath); err != nil {
		return err
	}
	return nil
}
