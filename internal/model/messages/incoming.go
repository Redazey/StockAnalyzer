package messages

import (
	"context"
	"strings"

	"go.uber.org/zap"

	"stockanalyzer/config"
	types "stockanalyzer/internal/model/bottypes"
	companyDB "stockanalyzer/internal/model/db"
	"stockanalyzer/pkg/cache"
	"stockanalyzer/pkg/logger"
)

// MessageSender Интерфейс для работы с сообщениями.
type MessageSender interface {
	SendMessage(text string, userID int64) (int, error)
	ShowKeyboardButtons(text string, buttons types.TgKbRowButtons, userID int64) error
	ShowInlineButtons(text string, buttons []types.TgRowButtons, userID int64) (int, error)
	EditInlineButtons(text string, msgID int, userID int64, buttons []types.TgRowButtons) error
}

// UserDataStorage Интерфейс для работы с хранилищем данных.
type Storage interface {
	GetCompanies(context.Context) ([]string, error)
	GetCompInfo(context.Context) ([]types.CompanyInfo, error)
}

// Model Модель бота (клиент, хранилище, последние команды пользователя)
type Model struct {
	ctx      context.Context
	tgClient MessageSender // Клиент
	storage  Storage       // Хранилище информации
	cfg      *config.Enviroment
}

// New Генерация сущности для хранения клиента ТГ и хранилища пользователей
func New(ctx context.Context, tgClient MessageSender, storage *companyDB.Storage, cfg *config.Enviroment) *Model {
	return &Model{
		ctx:      ctx,
		tgClient: tgClient,
		storage:  storage,
		cfg:      cfg,
	}
}

func (s *Model) GetCtx() context.Context {
	return s.ctx
}

func (s *Model) SetCtx(ctx context.Context) {
	s.ctx = ctx
}

// IncomingMessage Обработка входящего сообщения.
func (s *Model) IncomingMessage(msg types.Message) error {
	var companies []string
	cCompanies, err := cache.ReadCache("ctgShorts")
	if err != nil {
		return err
	}

	if cCompanies == "" {
		companies, err = s.storage.GetCompanies(s.ctx)
		if err != nil {
			return err
		}

		if err := cache.SaveCache("companies", strings.Join(companies, " ")); err != nil {
			return err
		}
	} else {
		companies = strings.Split(cCompanies, " ")
	}

	var compInfo []types.CompanyInfo
	var cCompInfo []types.CompanyInfo
	err = cache.ReadMapCache("compInfo", &compInfo)
	if err != nil {
		return err
	}

	if compInfo == nil {
		cCompInfo, err = s.storage.GetCompInfo(s.ctx)
		if err != nil {
			return err
		}

		if err := cache.SaveMapCache("compInfo", cCompInfo); err != nil {
			return err
		}
	} else {
		compInfo = cCompInfo
	}

	// Распознавание стандартных команд.
	if isNeedReturn, err := CheckBotCommands(s, msg, companies); err != nil || isNeedReturn {
		if err != nil {
			logger.Error("Error while CheckBotCommands: ", zap.Error(err))
		}

		return err
	}

	if isNeedReturn, err := CallbacksCommands(s, msg, companies); err != nil || isNeedReturn {
		if err != nil {
			logger.Error("Error while CallbacksCommands: ", zap.Error(err))
		}

		return err
	}

	// Отправка ответа по умолчанию.
	_, err = s.tgClient.SendMessage(TxtUnknownCommand, msg.UserID)
	return err
}

func GetCtgInfoFromName(Name string, compInfo []types.CompanyInfo) types.CompanyInfo {
	for _, ctg := range compInfo {
		if ctg.Name == Name {
			return ctg
		}
	}

	return types.CompanyInfo{}
}

func GetCtgInfoFromID(ID int64, compInfo []types.CompanyInfo) types.CompanyInfo {
	for _, ctg := range compInfo {
		if ctg.ID == ID {
			return ctg
		}
	}

	return types.CompanyInfo{}
}
