package db

import (
	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"github.com/nilorg/pkg/logger"
)

// DataBaseConfig ...
type DataBaseConfig struct {
	SnowflakeNode int64
	DBType        string
	MasterAddress string
	LogFlag       bool
	MaxOpen       int
	MaxIdle       int
	SlaveAddress  []string
}
type DataBase struct {
	master        *gorm.DB
	slaves        []*gorm.DB
	slaveIndex    int
	snowflakeNode *snowflake.Node
}

// NewDataBase ...
func NewDataBase(conf DataBaseConfig) *DataBase {
	node, err := snowflake.NewNode(conf.SnowflakeNode)
	if err != nil {
		logger.Fatalf("NewDataBase snowflake:%v", err)
	}

	master := newGorm(conf.DBType, conf.MasterAddress, conf.LogFlag, conf.MaxOpen, conf.MaxIdle)

	var slaves []*gorm.DB
	slaveAddressLen := len(conf.SlaveAddress)
	if slaveAddressLen == 0 {
		slaves = append(slaves, master)
	} else {
		for i := 0; i < slaveAddressLen; i++ {
			slaves = append(slaves, newGorm(conf.DBType, conf.SlaveAddress[i], conf.LogFlag, conf.MaxOpen, conf.MaxIdle))
		}
	}
	return &DataBase{
		snowflakeNode: node,
		master:        master,
		slaves:        slaves,
		slaveIndex:    0,
	}
}

// newGorm 创建...
func newGorm(dbType, address string, logFlag bool, maxOpen, maxIdle int) *gorm.DB {

	db, err := gorm.Open(dbType, address)
	if err != nil {
		logger.Fatalf("初始化 %s 连接失败: %s ", dbType, err)
	}
	err = db.DB().Ping()
	if err != nil {
		logger.Fatalf("Ping %s 连接失败: %s ", dbType, err)
	}
	db.LogMode(logFlag)

	db.DB().SetMaxOpenConns(maxOpen)
	db.DB().SetMaxIdleConns(maxIdle)
	return db
}

// Open 打开
func (db *DataBase) Close() {
	err := db.master.Close()
	if err != nil {
		logger.Errorln(err)
	}

	slaveLen := len(db.slaves)
	for i := 0; i < slaveLen; i++ {
		err = db.slaves[i].Close()
		if err != nil {
			logger.Errorln(err)
		}
	}
}

// Master 主
func (db *DataBase) Master() *gorm.DB {
	return db.master
}

// Slave 从
func (db *DataBase) Slave() *gorm.DB {
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

// NewSnowflakeID 雪花ID
func (db *DataBase) NewSnowflakeID() snowflake.ID {
	return db.snowflakeNode.Generate()
}
