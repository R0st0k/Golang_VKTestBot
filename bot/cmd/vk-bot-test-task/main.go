package main

import (
	"VKTestBot/internal/pkg/app"
)

func main() {

	// Создание объекта приложения
	a, err := app.New()
	if err != nil {
		panic(err)
	}

	// Запуск приложения
	err = a.Run()
	if err != nil {
		panic(err)
	}
}
