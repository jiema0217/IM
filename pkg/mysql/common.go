package mysql

import (
	"IMProject/util"
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"sync"
	"time"
)

var (
	mySqlPool = make(map[string]*sql.DB)
	once      sync.Once
)

// MySqlConfig mySql配置
// 大部分业务场景都是读多写少，数据库配置一般都采用一主多从
// 如果出现写多的情况，会采用 MongoDB, 或用 redis、MQ 先写入，在异步刷入到 MySQL
type MySqlConfig struct {
	Write           MySqlConn     `yaml:"write"`
	Read            []MySqlConn   `yaml:"read"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
}

// MySqlConn mysql连接信息
type MySqlConn struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
}

// InitMysql 初始化数据库
func InitMysql(mySqlConfig map[string]MySqlConfig) {
	once.Do(func() {
		var err error
		for key, val := range mySqlConfig {
			val.Write.Password, err = util.RsaDecrypt(val.Write.Password)
			if err != nil {
				panic(err)
			}
			for idx, read := range val.Read {
				val.Read[idx].Password, err = util.RsaDecrypt(read.Password)
				if err != nil {
					panic(err)
				}
			}
			mySqlConfig[key] = val
		}
		for dbName, cfg := range mySqlConfig {
			writeDsn := dsnToStr(cfg.Write)
			db, err := gorm.Open(mysql.Open(writeDsn), &gorm.Config{})
			if err != nil {
				panic(fmt.Sprintf("database [%s] init fail, err = [%v]", dbName, err))
			}
			readDsns := make([]string, 0, len(cfg.Read))
			for _, read := range cfg.Read {
				readDsns = append(readDsns, dsnToStr(read))
			}
			if len(readDsns) == 0 {
				readDsns = append(readDsns, writeDsn)
			}
			replicas := make([]gorm.Dialector, 0, len(readDsns))
			for _, dsn := range readDsns {
				replicas = append(replicas, mysql.Open(dsn))
			}
			db.Use(dbresolver.Register(dbresolver.Config{
				Replicas: replicas,
				Policy:   dbresolver.RandomPolicy{},
			}))
			sqlDB, err := db.DB()
			if err != nil {
				panic(fmt.Sprintf("database [%s] conn fail, err = [%v]", dbName, err))
			}
			sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
			sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
			sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
			sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
			mySqlPool[dbName] = sqlDB
		}
	})
}

// GetDbPool 获取数据库连接池
func GetDbPool(dbName string) *sql.DB {
	if db, ok := mySqlPool[dbName]; ok {
		return db
	} else {
		panic(fmt.Sprintf("database [%s] pool not exist", dbName))
	}
}

// dsnToStr 将数据库连接配置转换为dsn字符串
func dsnToStr(conn MySqlConn) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		conn.User,
		conn.Password,
		conn.Host,
		conn.Port,
		conn.Database,
		conn.Charset,
	)
}
