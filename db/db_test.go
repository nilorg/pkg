package db

import (
	"testing"

	"github.com/nilorg/pkg/logger"
)

func TestDataBase(t *testing.T) {
	conf := DataBaseConfig{
		DBType:        "postgres",
		MasterAddress: "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=test sslmode=disable",
		LogFlag:       true,
		MaxOpen:       100,
		MaxIdle:       10,
		SlaveAddress: []string{
			"host=127.0.0.2 port=5432 user=postgres password=postgres dbname=test sslmode=disable",
			"host=127.0.0.3 port=5432 user=postgres password=postgres dbname=test sslmode=disable",
			"host=127.0.0.4 port=5432 user=postgres password=postgres dbname=test sslmode=disable",
		},
	}
	logger.Init()
	db := NewDataBase(conf, logger.Default())
	db.Master()
	db.Slave()
}
