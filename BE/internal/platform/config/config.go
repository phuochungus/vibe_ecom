package config

import "os"

type ServiceConfig struct {
	ServiceName  string
	HTTPPort     string
	RedisAddr    string
	KafkaBrokers string
	RabbitMQURL  string
	MySQLDSN     string
}

func Load(serviceName string, defaultPort string) ServiceConfig {
	return ServiceConfig{
		ServiceName:  getEnv("SERVICE_NAME", serviceName),
		HTTPPort:     getEnv("HTTP_PORT", defaultPort),
		RedisAddr:    getEnv("REDIS_ADDR", "redis:6379"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "kafka:9092"),
		RabbitMQURL:  getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		MySQLDSN:     getEnv("MYSQL_DSN", "root:root@tcp(mysql:3306)/"+serviceName+"?parseTime=true"),
	}
}

func getEnv(key string, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	return v
}
