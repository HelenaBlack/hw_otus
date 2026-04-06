package config

import (
	"os"

	"gopkg.in/yaml.v2"
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
	Level string `yaml:"level"` // error, warn, info, debug
}

// StorageConf описывает тип используемого хранилища.
type StorageConf struct {
	Type string `yaml:"type"` // memory или sql
}

// ServerConf содержит параметры HTTP и GRPC серверов.
type ServerConf struct {
	Host     string `yaml:"host"`     // адрес
	Port     int    `yaml:"port"`     // HTTP порт
	GrpcPort int    `yaml:"grpcPort"` // GRPC порт
}

// DBConf содержит параметры подключения к базе данных.
type DBConf struct {
	Host     string `yaml:"host"`     // адрес БД
	Port     int    `yaml:"port"`     // порт БД
	User     string `yaml:"user"`     // пользователь
	Password string `yaml:"password"` // пароль
	DBName   string `yaml:"dbname"`   // имя базы
}

// RabbitMQConf содержит параметры подключения к RabbitMQ.
type RabbitMQConf struct {
	Host     string `yaml:"host"`     // адрес
	Port     int    `yaml:"port"`     // порт
	User     string `yaml:"user"`     // пользователь
	Password string `yaml:"password"` // пароль
	Queue    string `yaml:"queue"`    // имя очереди
}

// SchedulerConf содержит параметры планировщика.
type SchedulerConf struct {
	ScanInterval int `yaml:"scanInterval"` // интервал сканирования БД (сек)
}

// SenderConf содержит параметры рассыльщика.
type SenderConf struct {
	WorkerCount int `yaml:"workerCount"` // количество воркеров
}

// NewConfigFromFile читает и парсит YAML-конфиг из файла.
func NewConfigFromFile(path string) (Config, error) {
	var cfg Config
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer func() { _ = f.Close() }()
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
