package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var configEnv map[string]string

func LoadEnv() {
	var err error
	if os.Getenv("MODE") != "production" {
		configEnv, err = godotenv.Read(".env")
		if err != nil {
			return
		}
	}
}

func GetEnv(key, fallback string) string {
	if value, ok := configEnv[key]; ok {
		return value
	}
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func GetEnvAsInt(key string, fallback int) int {
	if value, ok := configEnv[key]; ok {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Invalid integer value for key %s: %v\n", key, err)
			return fallback
		}
		return intValue
	}

	if val := os.Getenv(key); val != "" {
		if intValue, err := strconv.Atoi(val); err == nil {
			return intValue
		}
		log.Printf("Invalid integer in os.Getenv for %s: %v\n", key, val)
	}
	return fallback
}

func GetEnvAsBool(key string, fallback bool) bool {
	valStr := strings.ToLower(GetEnv(key, ""))
	if valStr == "true" || valStr == "1" || valStr == "yes" {
		return true
	} else if valStr == "false" || valStr == "0" || valStr == "no" {
		return false
	}
	return fallback
}

func GetEnvAsFloat(key string, fallback float64) float64 {
	valStr := GetEnv(key, "")
	if val, err := strconv.ParseFloat(valStr, 64); err == nil {
		return val
	}
	return fallback
}
