package configs

import "os"

// Функция для подставления значения по умолчанию для переменных среды
func Getenv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		value = defaultValue
	}
	return value
}
