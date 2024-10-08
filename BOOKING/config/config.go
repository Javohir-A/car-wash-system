package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server      ServerConfig
		Database    DatabaseConfig
		Redis       RedisConfig
		JWT         JWTConfig
		RabbitMQ    RabbitMQConfig
		MongoConfig MongoConfig
		Auth        string
		Booking     string
	}
	JWTConfig struct {
		SecretKey string
	}

	ServerConfig struct {
		Host string
		Port string
	}
	DatabaseConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	RedisConfig struct {
		Host string
		Port string
	}
	RabbitMQConfig struct {
		RabbitMQ string
	}
	MongoConfig struct {
		User     string
		Password string
		Host     string
		Port     string
		DBname   string
	}
)

func (c *Config) Load() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	c.Server.Host = os.Getenv("SERVER_HOST")
	c.Server.Port = os.Getenv("SERVER_PORT")

	c.Database.Host = os.Getenv("DB_HOST")
	c.Database.Port = os.Getenv("DB_PORT")
	c.Database.User = os.Getenv("DB_USER")
	c.Database.Password = os.Getenv("DB_PASS")
	c.Database.DBName = os.Getenv("DB_NAME")

	c.Redis.Host = os.Getenv("REDIS_HOST")
	c.Redis.Port = os.Getenv("REDIS_PORT")

	c.JWT.SecretKey = os.Getenv("JWT_SECRET_KEY")

	c.RabbitMQ.RabbitMQ = os.Getenv("RABBITMQ_URL")

	c.MongoConfig.User = os.Getenv("MONGO_USER")
	c.MongoConfig.Host = os.Getenv("MONGO_HOST")
	c.MongoConfig.Password = os.Getenv("MONGO_PASS")
	c.MongoConfig.Port = os.Getenv("MONGO_PORT")
	c.MongoConfig.DBname = os.Getenv("MONGO_DBNAME")

	c.Auth = os.Getenv("AUTH")
	return nil
}

func NewConfig() *Config {
	return &Config{}
}
