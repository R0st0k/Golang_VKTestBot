package handler

import (
	"VKTestBot/internal/app/event"
	"VKTestBot/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Обработчик приветственного сообщения
func HelloMessageHandler(ctx context.Context, object event.MessageNewObject, tools *Tools) error {
	// Проверка на соответствие
	reStart := regexp.MustCompile("{\"command\":\"start\"}").Match([]byte(object.Message.Payload))
	reBack := regexp.MustCompile("{\"command\":\"back\"}").Match([]byte(object.Message.Payload))
	reEmpty := regexp.MustCompile("^$").Match([]byte(object.Message.Payload))
	if !(reStart || reBack || (reEmpty && (*tools.Status)[strconv.FormatInt(int64(object.Message.FromID), 10)] == "")) {
		return nil
	}
	// Пользователь не должен вводить данные, ждём
	(*tools.Status)[strconv.FormatInt(int64(object.Message.FromID), 10)] = "Wait"

	text := ""
	if reStart {
		text = "Приветствую! Я бот, который поможет тебя сохранить все пароли и легко получить их обратно! Что бы ты хотел сделать?"
	} else {
		text = "Чем я могу помочь?"
	}

	keyboard := event.MessagesKeyboard{
		OneTime: false,
		Buttons: [][]event.MessagesKeyboardButton{
			{
				{
					Color: "primary",
					Action: event.MessagesKeyboardButtonAction{
						Type:    "text",
						Label:   "Добавить Пароль",
						Payload: "{\"command\":\"add_password\"}",
					},
				},
			},
			{
				{
					Color: "primary",
					Action: event.MessagesKeyboardButtonAction{
						Type:    "text",
						Label:   "Найти Пароль",
						Payload: "{\"command\":\"find_one\"}",
					},
				},
			},
			{
				{
					Color: "primary",
					Action: event.MessagesKeyboardButtonAction{
						Type:    "text",
						Label:   "Показать Все Ресурсы",
						Payload: "{\"command\":\"show_resources\"}",
					},
				},
			},
			{
				{
					Color: "primary",
					Action: event.MessagesKeyboardButtonAction{
						Type:    "text",
						Label:   "Удалить Пароль",
						Payload: "{\"command\":\"delete_password\"}",
					},
				},
			},
			{
				{
					Color: "secondary",
					Action: event.MessagesKeyboardButtonAction{
						Type:    "text",
						Label:   "Информация об Авторе",
						Payload: "{\"command\":\"about\"}",
					},
				},
			},
		},
	}

	empJSON, err := json.Marshal(keyboard)
	if err != nil {
		log.Fatalf(err.Error())
	}

	v := url.Values{}
	v.Add("user_id", strconv.FormatInt(int64(object.Message.FromID), 10))
	v.Add("random_id", strconv.FormatInt(int64(rand.NormFloat64()), 10))
	v.Add("peer_id", strconv.FormatInt(int64(object.Message.PeerID), 10))
	v.Add("message", text)
	v.Add("keyboard", fmt.Sprintf("%s", empJSON))

	_, err = tools.SendRequest(v.Encode())
	if err != nil {
		return err
	}

	return nil
}

// Обработчик сообщения "Об Авторе"
func AboutMessageHandler(ctx context.Context, object event.MessageNewObject, tools *Tools) error {
	matched := regexp.MustCompile("{\"command\":\"about\"}").Match([]byte(object.Message.Payload))
	if !matched {
		return nil
	}
	(*tools.Status)[strconv.FormatInt(int64(object.Message.FromID), 10)] = "Wait"

	keyboard := event.MessagesKeyboard{
		OneTime: false,
		Buttons: [][]event.MessagesKeyboardButton{
			{
				{
					Action: event.MessagesKeyboardButtonAction{
						Type:    "open_link",
						Label:   "Ссылка на GitHub",
						Payload: "{\"command\":\"github\"}",
						Link:    "https://github.com/R0st0k",
					},
				},
			},
			{
				{
					Action: event.MessagesKeyboardButtonAction{
						Type:    "open_link",
						Label:   "Ссылка на Резюме",
						Payload: "{\"command\":\"cv\"}",
						Link:    "https://drive.google.com/file/d/1RptQbPG1x3lag0Yox7Pz3_uSHfskzA00/view?usp=share_link",
					},
				},
			},
			{
				{
					Color: "secondary",
					Action: event.MessagesKeyboardButtonAction{
						Type:    "text",
						Label:   "Назад",
						Payload: "{\"command\":\"back\"}",
					},
				},
			},
		},
	}

	empJSON, err := json.Marshal(keyboard)
	if err != nil {
		log.Fatalf(err.Error())
	}

	v := url.Values{}
	v.Add("user_id", strconv.FormatInt(int64(object.Message.FromID), 10))
	v.Add("random_id", strconv.FormatInt(int64(rand.NormFloat64()), 10))
	v.Add("peer_id", strconv.FormatInt(int64(object.Message.PeerID), 10))
	v.Add("message", "Меня зовут Низовцов Ростислав, я студент 4 курса СПбГЭТУ \"ЛЭТИ\" по направлению программная инженерия. Люблю слушать музыку и играть на барабанах")
	v.Add("keyboard", fmt.Sprintf("%s", empJSON))

	_, err = tools.SendRequest(v.Encode())
	if err != nil {
		return err
	}

	return nil
}

// Обработчик сообщения о добавлении пароля
func AddPasswordMessageHandler(ctx context.Context, object event.MessageNewObject, tools *Tools) error {
	matched := regexp.MustCompile("{\"command\":\"add_password\"}").Match([]byte(object.Message.Payload))
	if !matched {
		return nil
	}
	// Ждём от пользователя ввод информации для добавления
	(*tools.Status)[strconv.FormatInt(int64(object.Message.FromID), 10)] = "Add"

	keyboard := event.MessagesKeyboard{
		OneTime: false,
		Buttons: [][]event.MessagesKeyboardButton{
			{
				{
					Color: "primary",
					Action: event.MessagesKeyboardButtonAction{
						Type:    "text",
						Label:   "Добавить Следующий Пароль",
						Payload: "{\"command\":\"add_password\"}",
					},
				},
			},
			{
				{
					Color: "secondary",
					Action: event.MessagesKeyboardButtonAction{
						Type:    "text",
						Label:   "Назад",
						Payload: "{\"command\":\"back\"}",
					},
				},
			},
		},
	}

	empJSON, err := json.Marshal(keyboard)
	if err != nil {
		log.Fatalf(err.Error())
	}

	v := url.Values{}
	v.Add("user_id", strconv.FormatInt(int64(object.Message.FromID), 10))
	v.Add("random_id", strconv.FormatInt(int64(rand.NormFloat64()), 10))
	v.Add("peer_id", strconv.FormatInt(int64(object.Message.PeerID), 10))
	v.Add("message", fmt.Sprintf("Введите ресурс, логин и пароль через пробел, чтобы получилось следующее сообщение:\n<Ресурс> <Логин> <Пароль>"))
	v.Add("keyboard", fmt.Sprintf("%s", empJSON))

	_, err = tools.SendRequest(v.Encode())
	if err != nil {
		return err
	}

	return nil
}

// Обработчик сообщения об удалении пароля
func DeletePasswordMessageHandler(ctx context.Context, object event.MessageNewObject, tools *Tools) error {
	matched := regexp.MustCompile("{\"command\":\"delete_password\"}").Match([]byte(object.Message.Payload))
	if !matched {
		return nil
	}
	// Ждем информацию о пароле для удаления
	(*tools.Status)[strconv.FormatInt(int64(object.Message.FromID), 10)] = "Delete"

	keyboard := event.MessagesKeyboard{
		OneTime: false,
		Buttons: [][]event.MessagesKeyboardButton{
			{
				{
					Color: "primary",
					Action: event.MessagesKeyboardButtonAction{
						Type:    "text",
						Label:   "Удалить Следующий Пароль",
						Payload: "{\"command\":\"delete_password\"}",
					},
				},
			},
			{
				{
					Color: "secondary",
					Action: event.MessagesKeyboardButtonAction{
						Type:    "text",
						Label:   "Назад",
						Payload: "{\"command\":\"back\"}",
					},
				},
			},
		},
	}

	empJSON, err := json.Marshal(keyboard)
	if err != nil {
		log.Fatalf(err.Error())
	}

	v := url.Values{}
	v.Add("user_id", strconv.FormatInt(int64(object.Message.FromID), 10))
	v.Add("random_id", strconv.FormatInt(int64(rand.NormFloat64()), 10))
	v.Add("peer_id", strconv.FormatInt(int64(object.Message.PeerID), 10))
	v.Add("message", fmt.Sprintf("Введите название ресурса, для которого необходимо забыть пароль"))
	v.Add("keyboard", fmt.Sprintf("%s", empJSON))

	_, err = tools.SendRequest(v.Encode())
	if err != nil {
		return err
	}

	return nil
}

// Обработчик сообщения о поиске пароля по названию ресурса
func FindOnePasswordMessageHandler(ctx context.Context, object event.MessageNewObject, tools *Tools) error {
	matched := regexp.MustCompile("{\"command\":\"find_one\"}").Match([]byte(object.Message.Payload))
	if !matched {
		return nil
	}
	// Ждем от пользователя ввод названия ресурса
	(*tools.Status)[strconv.FormatInt(int64(object.Message.FromID), 10)] = "Find"

	keyboard := event.MessagesKeyboard{
		OneTime: false,
		Buttons: [][]event.MessagesKeyboardButton{
			{
				{
					Color: "primary",
					Action: event.MessagesKeyboardButtonAction{
						Type:    "text",
						Label:   "Найти Другой Пароль",
						Payload: "{\"command\":\"find_one\"}",
					},
				},
			},
			{
				{
					Color: "secondary",
					Action: event.MessagesKeyboardButtonAction{
						Type:    "text",
						Label:   "Назад",
						Payload: "{\"command\":\"back\"}",
					},
				},
			},
		},
	}

	empJSON, err := json.Marshal(keyboard)
	if err != nil {
		log.Fatalf(err.Error())
	}

	v := url.Values{}
	v.Add("user_id", strconv.FormatInt(int64(object.Message.FromID), 10))
	v.Add("random_id", strconv.FormatInt(int64(rand.NormFloat64()), 10))
	v.Add("peer_id", strconv.FormatInt(int64(object.Message.PeerID), 10))
	v.Add("message", fmt.Sprintf("Введите название ресурса, для которого необходимо найти пароль"))
	v.Add("keyboard", fmt.Sprintf("%s", empJSON))

	_, err = tools.SendRequest(v.Encode())
	if err != nil {
		return err
	}

	return nil
}

// Обработчик сообщения о демонстрации всех хранящихся паролей
func ShowResourcesMessageHandler(ctx context.Context, object event.MessageNewObject, tools *Tools) error {
	matched := regexp.MustCompile("{\"command\":\"show_resources(?P<page>[0-9]*)\"}").FindAllStringSubmatch(object.Message.Payload, 1)
	if len(matched) == 0 {
		return nil
	}
	(*tools.Status)[strconv.FormatInt(int64(object.Message.FromID), 10)] = "Wait"

	// Количество ресурсов на одной странице
	var n int64 = 5
	var page int64 = 1

	// Получение номера страницы для демонтсрации
	if matched[0][1] != "" {
		var err error
		page, err = strconv.ParseInt(matched[0][1], 10, 64)
		if err != nil {
			return err
		}
	}

	passwords, err := tools.Repository.AllUserPassword(strconv.FormatInt(int64(object.Message.FromID), 10))
	if err != nil {
		return err
	}

	text := "Вот список Ваших ресурсов:\n"

	for i := int((page - 1) * n); i < int(page*n) && i < len(passwords); i++ {
		text = fmt.Sprintf("%s%s\n", text, passwords[i].Resource)
	}

	prev := page - 1
	if prev == 0 {
		prev = 1
	}
	next := page + 1
	if page*n >= int64(len(passwords)) {
		next = page
	}

	keyboard := event.MessagesKeyboard{
		OneTime: false,
		Buttons: [][]event.MessagesKeyboardButton{
			{
				{
					Color: "primary",
					Action: event.MessagesKeyboardButtonAction{
						Type:    "text",
						Label:   "Предыдущая Страница",
						Payload: fmt.Sprintf("{\"command\":\"show_resources%d\"}", prev),
					},
				},
			},
			{
				{
					Color: "primary",
					Action: event.MessagesKeyboardButtonAction{
						Type:    "text",
						Label:   "Следующая Страница",
						Payload: fmt.Sprintf("{\"command\":\"show_resources%d\"}", next),
					},
				},
			},
			{
				{
					Color: "secondary",
					Action: event.MessagesKeyboardButtonAction{
						Type:    "text",
						Label:   "Назад",
						Payload: "{\"command\":\"back\"}",
					},
				},
			},
		},
	}

	empJSON, err := json.Marshal(keyboard)
	if err != nil {
		log.Fatalf(err.Error())
	}

	v := url.Values{}
	v.Add("user_id", strconv.FormatInt(int64(object.Message.FromID), 10))
	v.Add("random_id", strconv.FormatInt(int64(rand.NormFloat64()), 10))
	v.Add("peer_id", strconv.FormatInt(int64(object.Message.PeerID), 10))
	v.Add("message", text)
	v.Add("keyboard", fmt.Sprintf("%s", empJSON))

	_, err = tools.SendRequest(v.Encode())
	if err != nil {
		return err
	}

	return nil
}

// бработчик сообщений для манипуляции с данными
func GetInfoAboutPasswordMessageHandler(ctx context.Context, object event.MessageNewObject, tools *Tools) error {
	matched := regexp.MustCompile("^$").Match([]byte(object.Message.Payload))
	if !(matched && (*tools.Status)[strconv.FormatInt(int64(object.Message.FromID), 10)] != "") {
		return nil
	}

	// Режим поличения информации
	mod := (*tools.Status)[strconv.FormatInt(int64(object.Message.FromID), 10)]
	text := "Что-то пошло не так, попробуй еще раз!"
	// Метка для добавления сообщения как временного
	isNeedDelete := false

	(*tools.Status)[strconv.FormatInt(int64(object.Message.FromID), 10)] = "Wait"

	switch mod {
	case "Add":
		// Получение данных для добавления
		res := regexp.MustCompile("^(?P<res>[a-zA-Z0-9_!@#$%^&*]+) (?P<login>[a-zA-Z0-9_!@#$%^&*]+) (?P<pass>[a-zA-Z0-9_!@#$%^&*]+)$").
			FindAllStringSubmatch(object.Message.Text, -1)

		if len(res) != 1 {
			text = "В этом сообщении не найдено совпадений или найдено слишком много. Попробуй еще раз, нажав кнопку"
			break
		}

		err := tools.Repository.Create(repository.Password{
			UserID:   strconv.FormatInt(int64(object.Message.FromID), 10),
			Resource: strings.ToLower(res[0][1]),
			Login:    res[0][2],
			Password: res[0][3],
		})

		if err != nil {
			return err
		}

		text = "Пароль успешно сохранен! Не забудь удалить своё сообщение для сохранности"

	case "Find":
		pass, err := tools.Repository.GetByUserAndResource(strconv.FormatInt(int64(object.Message.FromID), 10),
			strings.ToLower(object.Message.Text))

		if err == repository.ErrNotExists {
			text = "Не удалось найти пароль от данного ресурса. Перепроверьте ввод"
			break
		} else {
			if err != nil {
				return err
			}
		}

		isNeedDelete = true
		text = fmt.Sprintf("Ресурс: %s\nЛогин: %s\nПароль: %s", object.Message.Text, pass.Login, pass.Password)

	case "Delete":
		err := tools.Repository.Delete(strconv.FormatInt(int64(object.Message.FromID), 10),
			strings.ToLower(object.Message.Text))

		if err == repository.ErrDeleteFailed {
			text = "Не удалось найти пароль от данного ресурса. Перепроверьте ввод"
			break
		} else {
			if err != nil {
				return err
			}
		}

		text = "Пароль успешно удален!"
	}

	v := url.Values{}
	v.Add("user_id", strconv.FormatInt(int64(object.Message.FromID), 10))
	v.Add("random_id", strconv.FormatInt(int64(rand.NormFloat64()), 10))
	v.Add("peer_id", strconv.FormatInt(int64(object.Message.PeerID), 10))
	v.Add("message", text)

	response, err := tools.SendRequest(v.Encode())
	if err != nil {
		return err
	}
	if isNeedDelete {
		tools.ToDeleteMessages = append(tools.ToDeleteMessages, ToDeleteMessage{
			Created:   time.Now(),
			MessageID: response.Response,
			Message:   text,
			PeerID:    object.Message.PeerID})
	}

	return nil
}
