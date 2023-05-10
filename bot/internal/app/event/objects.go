package event

import "encoding/json"

// Структура для получения сообщения
type MessageNewObject struct {
	Message    MessagesMessage `json:"message"`
	ClientInfo ClientInfo      `json:"client_info"`
}

type MessageDeleteObject struct {
	MessageIDs   string `json:"message_ids" url:"message_ids"`       // Перечисление ID Собщений для удаления, через запятую
	DeleteForAll int    `json:"delete_for_all" url:"delete_for_all"` // 1 для удаления у всех и 0 для удаления "у себя"
	PeerID       int    `json:"peer_id" url:"peer_id"`               // ID беседы
}

type MessageSendObject struct {
	UserID      int               `json:"user_id" url:"user_id"`     // ID Получателя сообщения
	RandomID    int               `json:"random_id" url:"random_id"` // Случайное число для контроля уникальности сообщения
	PeerID      int               `json:"peer_id" url:"peer_id"`     // ID беседы
	Message     string            `json:"message" url:"message"`
	Attachments []json.RawMessage `json:"attachments" url:"-"` // Вложения
	Keyboard    MessagesKeyboard  `json:"keyboard" url:"-"`    // Настройки клавиатуры ответов для пользователя
}

type ClientInfo struct {
	ButtonActions  []string `json:"button_actions"`
	Keyboard       bool     `json:"keyboard"`
	InlineKeyboard bool     `json:"inline_keyboard"`
	Carousel       bool     `json:"carousel"`
	LangID         int      `json:"lang_id"`
}

type MessagesMessage struct {
	Attachments           []json.RawMessage `json:"attachments"`             // Вложения
	ConversationMessageID int               `json:"conversation_message_id"` // ID сообщения в беседе
	Date                  int               `json:"date"`                    // Дата
	FromID                int               `json:"from_id"`                 // ID отправителя
	FwdMessages           []MessagesMessage `json:"fwd_Messages"`            // Структура, хранящая дерево пересланных сообщений
	ReplyMessage          *MessagesMessage  `json:"reply_message"`           // Сообщение, на которое "ответили"
	ID                    int               `json:"id"`                      // ID сообщения
	Important             bool              `json:"important"`               // Важное ли
	IsHidden              bool              `json:"is_hidden"`               // Скрыто ли
	Out                   int               `json:"out"`                     // Было ли отправлена нами
	Keyboard              MessagesKeyboard  `json:"keyboard"`                // Клавиатура ответов
	Payload               string            `json:"payload"`                 // Вшитая в сообщение информация
	PeerID                int               `json:"peer_id"`                 // ID беседы
	RandomID              int               `json:"random_id"`               // Число для контроля уникальности
	Text                  string            `json:"text"`                    // Текст сообщения
}

type MessagesKeyboard struct {
	AuthorID int                        `json:"author_id,omitempty" url:"author_id,omitempty"`
	Buttons  [][]MessagesKeyboardButton `json:"buttons" url:"buttons"`                   // Настройка кнопок
	OneTime  bool                       `json:"one_time" url:"one_time"`                 // Отключить ли после первого сообщения
	Inline   bool                       `json:"inline,omitempty" url:"inline,omitempty"` // РАсположение кнопок внутри сообщения
}

type MessagesKeyboardButton struct {
	Action MessagesKeyboardButtonAction `json:"action" url:"action"`                   // Действие кнопки
	Color  string                       `json:"color,omitempty" url:"color,omitempty"` // Цвет кнопки
}

type MessagesKeyboardButtonAction struct {
	AppID   int    `json:"app_id,omitempty" url:"app_id,omitempty"` // ID VK App
	Hash    string `json:"hash,omitempty" url:"hash,omitempty"`
	Label   string `json:"label,omitempty" url:"label,omitempty"` // Название кнопки
	OwnerID int    `json:"owner_id,omitempty" url:"owner_id,omitempty"`
	Payload string `json:"payload,omitempty" url:"payload,omitempty"` // Вшитая в кнопку информация
	Type    string `json:"type" url:"type"`                           // Тип ответа от кнопки
	Link    string `json:"link,omitempty" url:"link,omitempty"`       // Ссылка на внешний ресурс
}
