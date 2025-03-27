package companyDB

import (
	"context"
	"fmt"
	types "stockanalyzer/internal/model/bottypes"
	"stockanalyzer/pkg/cache"
	"stockanalyzer/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Storage - Тип для хранилища информации о компаниях.
type Storage struct {
	db *mongo.Database
}

// db - *mongo.Client - ссылка на подключение к БД.
func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		db: db,
	}
}

func (s Storage) GetCompanies(ctx context.Context) ([]string, error) {
	cCollection := s.db.Collection("Companies")
	var compInfo []types.CompanyInfo
	var companies []string

	// Passing nil as the filter matches all documents in the collection
	cur, err := cCollection.Find(ctx, nil)
	if err != nil {
		logger.Error("Ошибка при поиске в БД ", zap.Error(err))
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		var company types.CompanyInfo
		if err := cur.Decode(&company); err != nil {
			logger.Error("Ошибка при декодировании информации о компании ", zap.Error(err))
		}

		if err := cache.SaveMapCache(fmt.Sprintf("%v_info", company.Name), company); err != nil {
			logger.Error("Ошибка при сохранении в кэш", zap.Error(err))
		}

		compInfo = append(compInfo, company)
		companies = append(companies, company.Name)
	}

	if err := cur.Err(); err != nil {
		logger.Error("Ошибка при работе с БД", zap.Error(err))
	}

	// Close the cursor once finished
	cur.Close(ctx)

	if err := cache.SaveCache("companies", companies); err != nil {
		return nil, err
	}
	if err := cache.SaveMapCache("companies_info", compInfo); err != nil {
		logger.Error("Ошибка при сохранении в кэш", zap.Error(err))
	}

	return []string{"MTC", "Gazprom", "Yandex"}, nil
}

func (s Storage) GetCompInfo(ctx context.Context) ([]types.CompanyInfo, error) {
	return nil, nil
}

func (s Storage) AddCompany(ctx context.Context, compInfo types.CompanyInfo) (bool, error) {
	cCollection := s.db.Collection("Companies")
	res, err := cCollection.InsertOne(ctx, compInfo)
	if err != nil {
		return false, err
	}

	if err := cache.SaveCache(fmt.Sprintf("%v_compInfo", res), compInfo); err != nil {
		return true, err
	}

	return true, nil
}
