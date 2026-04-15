package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

// Config описывает структуру основного конфига приложения.
type Config struct {
	Logger    LoggerConf    `yaml:"logger"`    // параметры логирования
	Storage   StorageConf   `yaml:"storage"`   // параметры хранилища
	Server    ServerConf    `yaml:"server"`    // параметры HTTP-сервера
	DB        DBConf        `yaml:"db"`        // параметры БД
	RabbitMQ  RabbitMQConf  `yaml:"rabbitmq"`  // параметры RabbitMQ
	Scheduler SchedulerConf `yaml:"scheduler"` // параметры планировщика
	Sender    SenderConf    `yaml:"sender"`    // параметры рассыльщика
}

// LoggerConf содержит параметры логирования.
type LoggerConf struct {
	Level string `yaml:"level" env:"LOGGER_LEVEL" env-default:"info"` // error, warn, info, debug
}

// StorageConf описывает тип используемого хранилища.
type StorageConf struct {
	Type string `yaml:"type" env:"STORAGE_TYPE" env-default:"memory"` // memory или sql
}

// ServerConf содержит параметры HTTP и GRPC серверов.
type ServerConf struct {
	Host     string `yaml:"host" env:"SERVER_HOST" env-default:"0.0.0.0"`
	Port     int    `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
	GrpcPort int    `yaml:"grpcPort" env:"SERVER_GRPC_PORT" env-default:"50051"`
}

// DBConf содержит параметры подключения к базе данных.
type DBConf struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"DB_USER" env-default:"calendar"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	DBName   string `yaml:"dbname" env:"DB_NAME" env-default:"calendar"`
}

// RabbitMQConf содержит параметры подключения к RabbitMQ.
type RabbitMQConf struct {
	Host     string `yaml:"host" env:"RABBITMQ_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"RABBITMQ_PORT" env-default:"5672"`
	User     string `yaml:"user" env:"RABBITMQ_USER" env-default:"guest"`
	Password string `yaml:"password" env:"RABBITMQ_PASSWORD" env-default:"guest"`
	Queue    string `yaml:"queue" env:"RABBITMQ_QUEUE" env-default:"calendar"`
}

// SchedulerConf содержит параметры планировщика.
type SchedulerConf struct {
	ScanInterval int `yaml:"scanInterval" env:"SCHEDULER_SCAN_INTERVAL" env-default:"5"` // интервал сканирования БД (сек)
}

// SenderConf содержит параметры рассыльщика.
type SenderConf struct {
	WorkerCount int `yaml:"workerCount" env:"SENDER_WORKER_COUNT" env-default:"5"` // количество воркеров
}

// NewConfigFromFile читает и парсит YAML-конфиг из файла.
// Если путь к файлу пустой, читает только из переменных окружения.
func NewConfigFromFile(path string) (Config, error) {
	var cfg Config
	var err error
	if path != "" {
		err = cleanenv.ReadConfig(path, &cfg)
	} else {
		err = cleanenv.ReadEnv(&cfg)
	}
	return cfg, err
}
