package messages

import (
	"fmt"
	types "stockanalyzer/internal/model/bottypes"
	"stockanalyzer/pkg/cache"
)

// Распознавание стандартных команд бота.
func CheckBotCommands(s *Model, msg types.Message, companies []string) (bool, error) {
	switch msg.Text {
	case "/start":
		displayName := msg.UserDisplayName
		if len(displayName) == 0 {
			displayName = msg.UserName
		}

		if err := s.tgClient.ShowKeyboardButtons(fmt.Sprintf(TxtStart, displayName), BtnStart, msg.UserID); err != nil {
			return true, err
		}

		return true, nil
	case "Companies":
		var btns []types.TgRowButtons
		for _, ctg := range companies {
			btns = append(btns, types.TgRowButtons{types.TgInlineButton{DisplayName: ctg, Value: ctg}})
		}

		lastMsgID, err := s.tgClient.ShowInlineButtons(TxtCompanies, btns, msg.UserID)
		if err != nil {
			return true, err
		}

		if err := cache.SaveCache(fmt.Sprintf("%v_inlinekbMsg", msg.UserID), lastMsgID); err != nil {
			return true, err
		}

		return true, nil
	}

	// Команда не распознана.
	return false, nil
}
