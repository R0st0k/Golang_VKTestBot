package handler

import (
	"VKTestBot/internal/app/api"
	event2 "VKTestBot/internal/app/event"
	"VKTestBot/internal/app/parser"
	"VKTestBot/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

// Структура для хранения и предоставления инструментов для обработчиков
type HList struct {
	Tools      *Tools                                                         // Инструменты для отправки запросов и обработки временных сообщений
	MessageNew []func(context.Context, event2.MessageNewObject, *Tools) error // Список обработчиков на соответствующее событие
}

type Tools struct {
	API              api.API               // Доступ к API VK
	Status           *map[string]string    // Хранение текущего состояния пользователей
	Repository       repository.Repository // Доступ к хранилищу данных
	ToDeleteMessages []ToDeleteMessage     // Список временных сообщений
	ToDeleteTime     time.Duration         // Время жизни временных сообщений
}

type ToDeleteMessage struct {
	MessageID int
	Message   string
	PeerID    int
	Created   time.Time
}

// Отправка запросов о написании сообщения для обработчиков
func (tl *Tools) SendRequest(params string) (response parser.MessageResponse, err error) {
	resp, err := tl.API.VkAPICall("messages.send", params)
	if err != nil {
		return response, err
	}

	response, err = parser.ParseMessageResponse(resp.Body)
	if err != nil {
		return response, err
	}

	log.Printf("Bot send message with response: %d", response.Response)

	return response, nil
}

// Добавление обработчика в список
func (hl *HList) AddHandler(handlerType event2.EventType, f func(context.Context, event2.MessageNewObject, *Tools) error) {
	switch handlerType {
	case event2.EventMessageNew:
		if hl.MessageNew != nil {
			hl.MessageNew = append(hl.MessageNew, f)
		} else {
			hl.MessageNew = []func(context.Context, event2.MessageNewObject, *Tools) error{f}
		}
		log.Printf("EventMessageNew Handler was added")
	default:
		log.Printf("Can't add this handler (%s): unknown EventType", handlerType)
	}
	return
}

// Метод для распределения запросов по соответствующим обработчикам
func (hl *HList) Handle(ctx context.Context, e event2.GroupEvent) error {
	switch e.Type {
	case event2.EventMessageNew:
		var obj event2.MessageNewObject
		if err := json.Unmarshal(e.Object, &obj); err != nil {
			return err
		}

		for _, f := range hl.MessageNew {
			err := f(ctx, obj, hl.Tools)
			if err != nil {
				return err
			}
		}
	default:
		log.Printf("Can't handle this event group (%s): no handlers", e.Type)
	}

	return nil
}

// Проверка временных сообщений
func (hl *HList) CheckDeleteTimer() {
	for {
		if hl.Tools.ToDeleteMessages != nil {
			now := time.Now()
			count := 0
			for _, v := range hl.Tools.ToDeleteMessages {
				if now.After(v.Created.Add(hl.Tools.ToDeleteTime)) {
					count++
					err := hl.DeleteMessage(v)
					if err != nil {
						panic(err)
					}
				} else {
					err := hl.UpdateTimerOnDeletingMessages(v)
					if err != nil {
						panic(err)
					}
				}
			}
			hl.Tools.ToDeleteMessages = hl.Tools.ToDeleteMessages[count:]
			time.Sleep(2 * time.Second)
		}
	}
}

// Удаление временных сообщений
func (hl *HList) DeleteMessage(message ToDeleteMessage) error {
	v := url.Values{}
	v.Add("message_ids", strconv.FormatInt(int64(message.MessageID), 10))
	v.Add("delete_for_all", "1")
	v.Add("peer_id", strconv.FormatInt(int64(message.PeerID), 10))

	_, err := hl.Tools.API.VkAPICall("messages.delete", v.Encode())
	if err != nil {
		return err
	}

	return nil
}

// Обновление временных сообщений
func (hl *HList) UpdateTimerOnDeletingMessages(message ToDeleteMessage) error {
	res := regexp.MustCompile(`^.*(Сообщение удалится через [0-9]+ секунд\(ы\))$`).
		FindAllStringSubmatch(message.Message, 1)

	text := ""
	if len(res) == 0 {
		text = fmt.Sprintf("%s\n\nСообщение удалится через %d секунд(ы)", message.Message, message.Created.Add(hl.Tools.ToDeleteTime).Unix()-time.Now().Unix())
	} else {
		text = fmt.Sprintf("%sСообщение удалится через %d секунд(ы)", res[0][0], message.Created.Add(hl.Tools.ToDeleteTime).Unix()-time.Now().Unix())
	}

	v := url.Values{}
	v.Add("message_id", strconv.FormatInt(int64(message.MessageID), 10))
	v.Add("message", text)
	v.Add("peer_id", strconv.FormatInt(int64(message.PeerID), 10))

	_, err := hl.Tools.API.VkAPICall("messages.edit", v.Encode())
	if err != nil {
		return err
	}

	return nil
}
