package bottypes

type Empty struct{}

type CompanyInfo struct {
	ID        int64
	Name      string
	Percent   float64
	Advice    bool
	OtherInfo string
}

// Message Структура сообщения для обработки.
type Message struct {
	Text            string
	MessageID       int
	UserID          int64
	UserName        string
	UserDisplayName string
	IsCallback      bool
	CallbackMsgID   int
}

// Типы для описания состава кнопок телеграм сообщения.
// Кнопка сообщения.
type TgInlineButton struct {
	DisplayName string
	Value       string
	URL         string
}

// Строка с кнопками сообщения.
type TgRowButtons []TgInlineButton

// Типы для описания состава кнопок телеграм сообщения.
// Кнопка сообщения.
type TgKeyboardButton struct {
	Text string
}

// Строка с кнопками сообщения.
type TgKbRowButtons []TgKeyboardButton
