package companyDB

import (
	"context"
	types "stockanalyzer/internal/model/bottypes"

	"go.mongodb.org/mongo-driver/mongo"
)

// Storage - Тип для хранилища информации о компаниях.
type Storage struct {
	db *mongo.Client
}

// db - *mongo.Client - ссылка на подключение к БД.
func NewStorage(db *mongo.Client) *Storage {
	return &Storage{
		db: db,
	}
}

func (s Storage) GetCompanies(ctx context.Context) ([]string, error) {
	return []string{"MTC", "Gazprom", "Yandex"}, nil
}

func (s Storage) GetCompInfo(ctx context.Context) ([]types.CompanyInfo, error) {
	return nil, nil
}
