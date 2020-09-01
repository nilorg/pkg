package db

import (
	"log"

	gormV1 "github.com/jinzhu/gorm"
	nlog "github.com/nilorg/sdk/log"
)

// DataBaseConfig ...
type DataBaseConfig struct {
	DBType        string
	MasterAddress string
	LogFlag       bool
	MaxOpen       int
	MaxIdle       int
	SlaveAddress  []string
}

// DataBase ...
type DataBase struct {
	master     *gormV1.DB
	slaves     []*gormV1.DB
	slaveIndex int
	log        nlog.Logger
}

// NewDataBase ...
func NewDataBase(conf DataBaseConfig, log nlog.Logger) *DataBase {
	master := newGorm(conf.DBType, conf.MasterAddress, conf.LogFlag, conf.MaxOpen, conf.MaxIdle)

	var slaves []*gormV1.DB
	slaveAddressLen := len(conf.SlaveAddress)
	if slaveAddressLen == 0 {
		slaves = append(slaves, master)
	} else {
		for i := 0; i < slaveAddressLen; i++ {
			slaves = append(slaves, newGorm(conf.DBType, conf.SlaveAddress[i], conf.LogFlag, conf.MaxOpen, conf.MaxIdle))
		}
	}
	return &DataBase{
		master:     master,
		slaves:     slaves,
		slaveIndex: 0,
		log:        log,
	}
}

// newGorm 创建...
func newGorm(dbType, address string, logFlag bool, maxOpen, maxIdle int) *gormV1.DB {

	db, err := gormV1.Open(dbType, address)
	if err != nil {
		log.Fatalf("初始化 %s 连接失败: %s ", dbType, err)
	}
	err = db.DB().Ping()
	if err != nil {
		log.Fatalf("Ping %s 连接失败: %s ", dbType, err)
	}
	db.LogMode(logFlag)

	db.DB().SetMaxOpenConns(maxOpen)
	db.DB().SetMaxIdleConns(maxIdle)
	return db
}

// Close 关闭
func (db *DataBase) Close() {
	err := db.master.Close()
	if err != nil {
		db.log.Errorln(err)
	}

	slaveLen := len(db.slaves)
	for i := 0; i < slaveLen; i++ {
		err = db.slaves[i].Close()
		if err != nil {
			db.log.Errorln(err)
		}
	}
}

// Master 主
func (db *DataBase) Master() *gormV1.DB {
	return db.master
}

// Slave 从
func (db *DataBase) Slave() *gormV1.DB {
	slaveLen := len(db.slaves)
	if slaveLen == 0 {
		return db.Master()
	} else if slaveLen == 1 {
		return db.slaves[0]
	} else {
		slave := db.slaves[db.slaveIndex]

		if db.slaveIndex+1 >= slaveLen {
			db.slaveIndex = 0
		} else {
			db.slaveIndex++
		}
		return slave
	}
}
