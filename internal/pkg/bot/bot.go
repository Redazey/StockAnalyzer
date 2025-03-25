package bot

import (
	"context"
	"stockanalyzer/config"
	"stockanalyzer/internal/clients/tg"
	companyDB "stockanalyzer/internal/model/db"
	"stockanalyzer/internal/model/messages"
	"stockanalyzer/pkg/cache"
	"stockanalyzer/pkg/logger"
	db "stockanalyzer/pkg/mongo"
)

type App struct {
	tgClient *tg.Client
	storage  *companyDB.Storage
	msgModel *messages.Model
}

func Init() (*App, error) {
	a := &App{}

	ctx := context.Background()

	cfg, err := config.NewEnv()
	if err != nil {
		return nil, err
	}

	logger.Init(cfg.LoggerLevel, "")

	err = db.Init(cfg.DB.DBAddr, cfg.DB.DBPort)
	if err != nil {
		return nil, err
	}

	err = cache.Init(cfg.Redis.RedisAddr+":"+cfg.Redis.RedisPort, cfg.Redis.RedisPassword, 0, cfg.Redis.EXTime)
	if err != nil {
		return nil, err
	}

	a.tgClient, err = tg.New(cfg.TgToken, tg.HandlerFunc(tg.ProcessingMessages))
	if err != nil {
		return nil, err
	}

	a.storage = companyDB.NewStorage(db.GetDBConn())
	a.msgModel = messages.New(ctx, a.tgClient, a.storage, cfg)

	return a, nil
}

func (a *App) Run() error {
	logger.Info("Запуск бота")

	a.tgClient.ListenUpdates(a.msgModel)

	return nil
}
