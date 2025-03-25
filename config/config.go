package config

import (
	"log"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Файл переменных окружения
type Enviroment struct {
	LoggerLevel string `env:"loggerMode" envDefault:"debug"`
	TgToken     string `env:"TG_TOKEN,required"`
	DB          DB
	Redis       Redis
}

type DB struct {
	DBAddr string `env:"MONGO_ADDR,required"`
	DBPort string `env:"MONGO_PORT,required"`
}

type Redis struct {
	RedisAddr     string        `env:"REDIS_ADDR,required"`
	RedisPort     string        `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword string        `env:"REDIS_PASSWORD,required"`
	RedisDBId     int           `env:"REDIS_DB_ID,required"`
	EXTime        time.Duration `env:"CACHE_EX_TIME" envDefault:"15m"`
}

var enviroment Enviroment

/*
Структура env файла

	-------GENERAL------
	LoggerLevel   string
	TgToken       string
	---------DB---------
	DBPort        string
	DBAddr        string
	-------REDIS--------
	RedisAddr     string
	RedisPort     string
	RedisPassword string
	RedisDBid     int
	CacheEXTime   int
*/
func NewEnv(envPath ...string) (*Enviroment, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Файл .env не найден: %s", err)
	}

	err = env.Parse(&enviroment)
	if err != nil {
		return nil, err
	}
	err = env.Parse(&enviroment.Redis)
	if err != nil {
		return nil, err
	}
	err = env.Parse(&enviroment.DB)
	if err != nil {
		return nil, err
	}

	return &enviroment, nil
}

func GetEnv() *Enviroment {
	return &Enviroment{}
}
