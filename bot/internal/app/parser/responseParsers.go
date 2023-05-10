package parser

import (
	"VKTestBot/internal/app/event"
	"encoding/json"
	"errors"
	"io"
	"log"
	"strconv"
)

// Структура ответа на сообщения бота
type MessageResponse struct {
	Response int `json:"response"`
}

// Структура ответа на запрос информации о сервере
type ServerInfo struct {
	Key    string
	Server string
	Ts     string
}

// Структура ответа сервера событий
type Response struct {
	Ts      string             `json:"ts"`
	Updates []event.GroupEvent `json:"updates"`
	Failed  int                `json:"failed"`
}

func ParseMessageResponse(reader io.Reader) (response MessageResponse, err error) {
	decoder := json.NewDecoder(reader)
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return response, err
		}

		t, ok := token.(string)
		if !ok {
			continue
		}

		switch t {
		case "response":
			raw, err := decoder.Token()
			if err != nil {
				return response, err
			}

			response.Response = int(raw.(float64))
		}
	}

	return response, err
}

func ParseEventResponse(reader io.Reader) (response Response, err error) {
	decoder := json.NewDecoder(reader)
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return response, err
		}

		t, ok := token.(string)
		if !ok {
			continue
		}

		switch t {
		case "failed":
			raw, err := decoder.Token()
			if err != nil {
				return response, err
			}

			response.Failed = int(raw.(float64))
		case "updates":
			var updates []event.GroupEvent

			err = decoder.Decode(&updates)
			if err != nil {
				return response, err
			}

			response.Updates = updates
		case "ts":
			// can be a number in the response with "failed" field: {"ts":8,"failed":1}
			// or string, e.g. {"ts":"8","updates":[]}
			rawTs, err := decoder.Token()
			if err != nil {
				return response, err
			}

			if ts, isNumber := rawTs.(float64); isNumber {
				response.Ts = strconv.Itoa(int(ts))
			} else {
				response.Ts = rawTs.(string)
			}
		}
	}

	return response, err
}

func ParseServerInfoResponse(reader io.Reader) (response ServerInfo, err error) {
	decoder := json.NewDecoder(reader)
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return response, err
		}

		t, ok := token.(string)
		if !ok {
			continue
		}

		switch t {
		case "key":
			raw, err := decoder.Token()
			if err != nil {
				return response, err
			}

			response.Key = raw.(string)
		case "server":
			raw, err := decoder.Token()
			if err != nil {
				return response, err
			}

			response.Server = raw.(string)
		case "ts":
			raw, err := decoder.Token()
			if err != nil {
				return response, err
			}

			response.Ts = raw.(string)
		case "error_code":
			raw, err := decoder.Token()
			if err != nil {
				return response, err
			}

			log.Fatalf("Error in parseServerInfoResponse with code :%f\n", raw.(float64))
		case "error_msg":
			raw, err := decoder.Token()
			if err != nil {
				return response, err
			}

			log.Fatalf("Error in parseServerInfoResponse with message :%s\n", raw.(string))
		}

	}

	return response, err
}
