package longpoll

import (
	"VKTestBot/internal/app/api"
	"VKTestBot/internal/app/handler"
	"VKTestBot/internal/app/parser"
	"VKTestBot/repository"
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

// Структура для работы с VK LongPoll API
type LongPoll struct {
	Server string `url:"-"` // Название сервера для обращения
	Key    string `url:"-"` // Ключ для обращения
	Ts     string `url:"-"` // Номер последнего события
	Wait   int    `url:"-"` // Время обновления запроса
	API    api.API
	cancel context.CancelFunc `url:"-"`

	*handler.HList
}

func NewLongPoll() (*LongPoll, error) {
	lp := &LongPoll{
		API: api.API{
			AccessToken: os.Getenv("ACCESS_TOKEN"),
			GroupID:     os.Getenv("GROUP_ID"),
			V:           os.Getenv("V"),
			Client:      http.DefaultClient,
		},
		Wait: 25,
	}

	repository.Init()
	log.Printf("Repository initiated")

	lp.HList = new(handler.HList)
	lp.HList.Tools = new(handler.Tools)
	lp.HList.Tools.API = lp.API
	status := make(map[string]string, 10)
	lp.HList.Tools.Status = &status
	lp.HList.Tools.Repository = repository.NewSQLiteRepository()
	lp.HList.Tools.ToDeleteMessages = make([]handler.ToDeleteMessage, 0, 10)
	lp.HList.Tools.ToDeleteTime = 30 * time.Second

	go lp.HList.CheckDeleteTimer()

	err := lp.getLongPollServer()

	return lp, err
}

// Метод для получения данных о сервере для работы с ним
func (lp *LongPoll) getLongPollServer() error {

	resp, err := lp.API.VkAPICall("groups.getLongPollServer", "")

	response, err := parser.ParseServerInfoResponse(resp.Body)
	if err != nil {
		return err
	}

	lp.Server = response.Server
	lp.Key = response.Key
	lp.Ts = response.Ts

	log.Printf("ServerInfo was recieved")

	return nil
}

// Метод для проверки наличия необработанных событий
func (lp *LongPoll) checkEvent(ctx context.Context) (response parser.Response, err error) {
	u := fmt.Sprintf("%s?act=a_check&key=%s&ts=%s&wait=%d", lp.Server, lp.Key, lp.Ts, lp.Wait)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return response, err
	}

	resp, err := lp.API.Client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	response, err = parser.ParseEventResponse(resp.Body)
	if err != nil {
		return response, err
	}

	err = lp.checkEventResponse(response)

	return response, err
}

// Метод для проверки кода ошибки ответа сервера событий
func (lp *LongPoll) checkEventResponse(response parser.Response) (err error) {
	switch response.Failed {
	case 0:
		lp.Ts = response.Ts
	case 1:
		lp.Ts = response.Ts
	case 2:
		err = lp.getLongPollServer()
	case 3:
		err = lp.getLongPollServer()
	default:
		log.Printf("Response failed")
	}

	return
}

// Метод для запуска проверки событий в постоянном режиме
func (lp *LongPoll) Run(ctx context.Context) error {
	ctx, lp.cancel = context.WithCancel(ctx)

	for {
		select {
		case _, ok := <-ctx.Done():
			if !ok {
				return nil
			}
		default:
			resp, err := lp.checkEvent(ctx)
			if err != nil {
				return err
			}

			ctx = context.WithValue(ctx, "ts", resp.Ts)

			for _, event := range resp.Updates {
				err = lp.HList.Handle(ctx, event)
				if err != nil {
					return err
				}
			}
		}
	}
}

func (lp *LongPoll) Shutdown() {
	if lp.cancel != nil {
		lp.cancel()
	}
}
