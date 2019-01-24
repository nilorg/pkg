package db

import "testing"

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
	db := NewDataBase(conf)
	db.Master()
	db.Slave()
}
