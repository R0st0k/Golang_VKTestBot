package api

import (
	"fmt"
	"log"
	"net/http"
)

// Структура взаимодействия с VK API
type API struct {
	GroupID     string       `url:"group_id"`
	AccessToken string       `url:"access_token"`
	V           string       `url:"v"`
	Client      *http.Client `url:"-"`
}

// Метод для отправки запроса на сервер
func (api *API) VkAPICall(method string, params string) (response *http.Response, err error) {
	u := fmt.Sprintf("https://api.vk.com/method/%s?access_token=%s&group_id=%s&v=%s&%s",
		method, api.AccessToken, api.GroupID, api.V, params)

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return response, err
	}

	resp, err := api.Client.Do(req)
	if err != nil {
		return resp, err
	}

	log.Printf("VkAPICall was used with method: %s", method)

	return resp, err
}
