package messages

import (
	types "stockanalyzer/internal/model/bottypes"
)

// Команды стартовых действий.
var BtnStart = types.TgKbRowButtons{
	types.TgKeyboardButton{Text: "Companies"},
}

const (
	WorkersChatID = 00000000000 // -- настроить
	TxtStart      = `Привет, %v 👋. Это бот, который анализирует экономическое состояние компаний и прогнозирует рост или падение акций. 
	Внимание! Эти прогнозы не являются личной рекомендацией к инвистициям.`
	TxtCompanies      = "📰 Выбери компанию:"
	TxtUnknownCommand = "Неизвестная компания"
)
