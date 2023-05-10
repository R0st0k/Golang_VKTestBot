package app

import (
	"VKTestBot/internal/app/event"
	"VKTestBot/internal/app/handler"
	"VKTestBot/internal/app/longpoll"
	"context"
	"log"
)

type App struct {
	LongPoll *longpoll.LongPoll
}

func New() (*App, error) {
	lp, err := longpoll.NewLongPoll()
	if err != nil {
		return nil, err
	}

	// Добавление обработчиков
	lp.HList.AddHandler(event.EventMessageNew, handler.HelloMessageHandler)
	lp.HList.AddHandler(event.EventMessageNew, handler.AddPasswordMessageHandler)
	lp.HList.AddHandler(event.EventMessageNew, handler.FindOnePasswordMessageHandler)
	lp.HList.AddHandler(event.EventMessageNew, handler.ShowResourcesMessageHandler)
	lp.HList.AddHandler(event.EventMessageNew, handler.DeletePasswordMessageHandler)
	lp.HList.AddHandler(event.EventMessageNew, handler.AboutMessageHandler)
	lp.HList.AddHandler(event.EventMessageNew, handler.GetInfoAboutPasswordMessageHandler)

	log.Printf("Handlers added")

	return &App{LongPoll: lp}, nil
}

func (a *App) Run() error {
	log.Printf("Server Running")
	return a.LongPoll.Run(context.Background())
}

func (a *App) Close() error {
	a.LongPoll.Shutdown()
	return nil
}
