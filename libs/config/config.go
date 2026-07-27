package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Srv     *Server
	MySql   *MySql
	Kafka   *Kafka
	Service *Service
	JWT     *JWT
}

type Server struct {
	Name string
	Host string
	Port string
}

type JWT struct {
	Secret string
}

type Service struct {
	Auth         string
	User         string
	Order        string
	Payment      string
	Notification string
}

type MySql struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type Kafka struct {
	Broker  string
	Topic   string
	GroupID string
}

// Load reads configuration from the environment.
// Optional envFiles are loaded first (e.g. build/.env/.<app>.env); then a root .env if present.
// Already-set process env vars always win over file values (including Docker Compose env_file).
func Load(envFiles ...string) *Config {
	if len(envFiles) > 0 {
		_ = godotenv.Load(envFiles...)
	}
	_ = godotenv.Load()

	cfg := &Config{
		Srv: &Server{
			Name: firstNonEmpty(os.Getenv("NAME"), os.Getenv("APP_NAME"), "unknown"),
			Host: getEnv("HOST", "localhost"),
			Port: getEnv("PORT", "8080"),
		},
		MySql: &MySql{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			Database: getEnv("DB_NAME", "userdb"),
		},
		Kafka: &Kafka{
			Broker:  getEnv("KAFKA_BROKER", "localhost:9092"),
			Topic:   getEnv("KAFKA_TOPIC", "payment-events"),
			GroupID: getEnv("KAFKA_GROUP_ID", "notification-service"),
		},
		Service: &Service{
			Auth:  getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
			User:  getEnv("USER_SERVICE_URL", "http://localhost:8082"),
			Order: getEnv("ORDER_SERVICE_URL", "http://localhost:8083"),
			Payment: firstNonEmpty(
				os.Getenv("PAYMENT_SERVICE_URL"),
				os.Getenv("PAYMENT_SERVICE_ADDR"),
				"localhost:50051",
			),
			Notification: getEnv("NOTIFICATION_SERVICE_URL", ""),
		},
		JWT: &JWT{
			Secret: getEnv("JWT_SECRET", "some-random-jwt-key"),
		},
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// EnvFileForApp returns the path to the app env file under build/.env/.
// Layout: build/.env/.<app>.env (e.g. build/.env/.gateway.env).
func EnvFileForApp(app string) string {
	return fmt.Sprintf("build/.env/.%s.env", app)
}

func LoadForApp(app string) *Config {
	return Load(EnvFileForApp(app))
}
