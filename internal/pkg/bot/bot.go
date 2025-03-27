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
	"time"
)

type App struct {
	tgClient *tg.Client
	storage  *companyDB.Storage
	msgModel *messages.Model
}

func Init() (*App, error) {
	a := &App{}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	cfg, err := config.NewEnv()
	if err != nil {
		return nil, err
	}

	logger.Init(cfg.LoggerLevel, "")

	err = db.Init(ctx, cfg.DB.DBAddr, cfg.DB.DBPort, cfg.DB.DBName)
	if err != nil {
		return nil, err
	}

	err = cache.Init(ctx, cfg.Redis.RedisAddr+":"+cfg.Redis.RedisPort, cfg.Redis.RedisPassword, 0, cfg.Redis.EXTime)
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
