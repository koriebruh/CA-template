package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

// FOR LOAD ALL CONFIG

type Config struct {
	Server
	DataBase
	RedisDB
}

type Server struct {
	Host string
	Port string
}

type DataBase struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

type RedisDB struct {
	Addr     string
	Password string
	DB       int
	Protocol int
}

// <-- CONSTRUCTOR --> //

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error in load .env : ", err.Error())
	}

	// Convert the RedisDB.DB and RedisDB.Protocol to integers
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Error converting REDIS_DB to int: %v", err)
	}

	redisProtocol, err := strconv.Atoi(os.Getenv("REDIS_PROTOCOL"))
	if err != nil {
		log.Fatalf("Error converting REDIS_PROTOCOL to int: %v", err)
	}

	return &Config{
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		DataBase: DataBase{
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			Name: os.Getenv("DB_NAME"),
		},
		RedisDB: RedisDB{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASS"),
			DB:       redisDB,
			Protocol: redisProtocol,
		},
	}

}
