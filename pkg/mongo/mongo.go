package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Conn *mongo.Client

func Init(Address string, Port string) error {
	ctx := context.Background()
	connStr := fmt.Sprintf(`mongodb://%v@mongo:%v/`, Address, Port)
	Conn, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	if err != nil {
		return err
	}

	// Проверка подключения к базе данных
	err = Conn.Ping(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

/* принимает map - значения, которые нужно внести в БД и string - таблицу, в которую будем вносить значения
func PullData(table string, data map[string]map[string]interface{}) error {
	for _, keyData := range data {
		var (
			columns []string
			values  []string
		)

		for key, value := range keyData {
			columns = append(columns, key)
			values = append(values, fmt.Sprintf("%s", value))
		}
		cmdStr := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (?)`, table, strings.Join(columns, ", "))
		query, args, err := sqlx.In(cmdStr, values)

		if err != nil {
			return err
		}

		query = Conn.Rebind(query)
		_, err = Conn.Query(query, args...)
		if err != nil {
			return err
		}
	}

	return nil
}
*/

func GetDBConn() *mongo.Client {
	return Conn
}
