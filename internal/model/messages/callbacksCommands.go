package messages

import (
	"fmt"
	types "stockanalyzer/internal/model/bottypes"
	"stockanalyzer/pkg/cache"
	"strconv"

	"github.com/opentracing/opentracing-go"
)

// callbacks
func CallbacksCommands(s *Model, msg types.Message, companies []string) (bool, error) {
	cacheInlinekbMsg, err := cache.ReadCache(fmt.Sprintf("%v_inlinekbMsg", msg.UserID))
	if err != nil {
		return true, err
	}

	lastInlinekbMsg, err := strconv.Atoi(cacheInlinekbMsg)
	if err != nil {
		return true, err
	}

	if msg.IsCallback {
		span, ctx := opentracing.StartSpanFromContext(s.ctx, "callbacksCommands")
		s.ctx = ctx
		defer span.Finish()

		// Дерево callbacks начинающихся с Categories
		switch msg.Text {
		case "backToCtg":
			var btns []types.TgRowButtons
			//var err error

			for _, company := range companies {
				btns = append(btns, types.TgRowButtons{types.TgInlineButton{DisplayName: company, Value: company}})
			}

			if lastInlinekbMsg == 0 {
				lastMsgID, err := s.tgClient.ShowInlineButtons(TxtCompanies, btns, msg.UserID)
				if err != nil {
					return true, err
				}

				if err := cache.SaveCache(fmt.Sprintf("%v_inlinekbMsg", msg.UserID), lastMsgID); err != nil {
					return true, err
				}
			}

			return true, s.tgClient.EditInlineButtons(
				TxtCompanies,
				lastInlinekbMsg,
				msg.UserID,
				btns,
			)
		}
	}

	// Команда не опознана
	return false, nil
}
