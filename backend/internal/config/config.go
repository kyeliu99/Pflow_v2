package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP      HTTPConfig
	Database  DatabaseConfig
	Queue     QueueConfig
	Camunda   CamundaConfig
	Telemetry TelemetryConfig
}

type HTTPConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	DSN                   string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
}

type QueueConfig struct {
	URL         string
	Exchange    string
	RoutingKey  string
	ContentType string
}

type CamundaConfig struct {
	BaseURL  string
	Username string
	Password string
}

type TelemetryConfig struct {
	ServiceName string
}

func Load(path string) (Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	if path != "" {
		v.AddConfigPath(path)
	}
	v.AddConfigPath(".")
	v.SetEnvPrefix("PFLOW")
	v.AutomaticEnv()

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, fmt.Errorf("config: read file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("config: unmarshal: %w", err)
	}
	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("http.host", "0.0.0.0")
	v.SetDefault("http.port", 8080)
	v.SetDefault("http.readTimeout", "10s")
	v.SetDefault("http.writeTimeout", "10s")

	v.SetDefault("database.dsn", "postgres://pflow:pflow@localhost:5432/pflow?sslmode=disable")
	v.SetDefault("database.maxOpenConnections", 20)
	v.SetDefault("database.maxIdleConnections", 5)
	v.SetDefault("database.connectionMaxLifetime", "30m")

	v.SetDefault("queue.url", "amqp://guest:guest@localhost:5672/")
	v.SetDefault("queue.exchange", "pflow.workorders")
	v.SetDefault("queue.routingKey", "events")
	v.SetDefault("queue.contentType", "application/json")

	v.SetDefault("camunda.baseURL", "http://localhost:8081/engine-rest")
	v.SetDefault("camunda.username", "demo")
	v.SetDefault("camunda.password", "demo")

	v.SetDefault("telemetry.serviceName", "pflow-backend")
}
