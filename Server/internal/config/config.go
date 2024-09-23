package config

import (
	"os"
)

type Config struct {
    ServerAddress string
    ServerPort    string
    MongoURI      string
    // TODO: Add more configuration fields as needed
}

func Load() (*Config, error) {
    return &Config{
        ServerAddress: getEnv("SERVER_ADDRESS", "localhost"),
        ServerPort:    getEnv("SERVER_PORT", "8080"),
        MongoURI:      getEnv("MONGO_URI", "mongodb+srv://alex:alex@cluster0.t5mgz.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"),
    }, nil
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}