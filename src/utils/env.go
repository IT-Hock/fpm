package utils

import (
	"os"
	"strconv"
)

func GetEnv(key string, def string) string {
	envVar := os.Getenv(key)
	if envVar == "" {
		return def
	}

	return envVar
}

func GetEnvBool(key string, def bool) bool {
	envVar := os.Getenv(key)
	if envVar == "" {
		return def
	}

	return envVar == "true" || envVar == "1"
}

func GetEnvInt(key string, def int) int {
	envVar := os.Getenv(key)
	if envVar == "" {
		return def
	}

	parsedInt, err := strconv.ParseInt(envVar, 10, 32)
	if err != nil {
		return def
	}

	return int(parsedInt)
}

func GetEnvInt64(key string, def int64) int64 {
	envVar := os.Getenv(key)
	if envVar == "" {
		return def
	}

	parsedInt, err := strconv.ParseInt(envVar, 10, 64)
	if err != nil {
		return def
	}

	return parsedInt
}

func GetEnvFloat(key string, def float64) float64 {
	envVar := os.Getenv(key)
	if envVar == "" {
		return def
	}

	parsedFloat, err := strconv.ParseFloat(envVar, 64)
	if err != nil {
		return def
	}

	return parsedFloat
}

func GetEnvFloat32(key string, def float32) float32 {
	envVar := os.Getenv(key)
	if envVar == "" {
		return def
	}

	parsedFloat, err := strconv.ParseFloat(envVar, 32)
	if err != nil {
		return def
	}

	return float32(parsedFloat)
}
