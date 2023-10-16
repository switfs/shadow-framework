package datasourcemanager

import (
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/switfs/shadow-framework/logger"
	"github.com/switfs/shadow-framework/orm/datasource"
	"github.com/switfs/shadow-framework/orm/jinzhu/gorm"
	"github.com/switfs/shadow-framework/server"

	_ "github.com/switfs/shadow-framework/orm/jinzhu/gorm/dialects/mysql"
)

type TSQLLogger struct{}

func (slog TSQLLogger) Print(values ...interface{}) {
	vals := gorm.LogFormatter(values...)
	logger.Log.SqlDebug(vals...)
}

var dbManger datasource.IDatasourceManager
var shardingManager datasource.IShardingDatasourceManager

type TGormDataSourceManager struct {
	configs    []datasource.IDataSourceConfigure
	masterDB   *gorm.DB
	slaveDB    *gorm.DB
	shardingDB map[string]*gorm.DB
	Models     []interface{}
}

func newGormDataSourceManager() datasource.IDatasourceManager {
	if dbManger == nil {
		l.Lock()
		defer l.Unlock()
		if dbManger == nil {
			dbManger = &TGormDataSourceManager{}
		}
	}
	return dbManger
}

func newGormShardingDatasourceManager() datasource.IShardingDatasourceManager {
	if shardingManager == nil {
		l.Lock()
		defer l.Unlock()
		if shardingManager == nil {
			shardingManager = &TGormDataSourceManager{}
		}
	}
	return shardingManager
}

func (manager *TGormDataSourceManager) Configs(configs []datasource.IDataSourceConfigure) {
	manager.configs = configs
}

// auto create/migrate table if you want
func (manager *TGormDataSourceManager) RegisterModels(models ...interface{}) {
	manager.Models = append(manager.Models, models...)
}

func (manager *TGormDataSourceManager) Datasource() *gorm.DB {
	return manager.Master()
}

func (manager *TGormDataSourceManager) SDatasource(key string) *gorm.DB {
	if manager.shardingDB == nil {
		l.Lock()
		defer l.Unlock()
		if manager.shardingDB == nil {
			m := make(map[string]*gorm.DB)
			if manager.configs == nil {
				panic("configs is nil")
			}
			for _, config := range manager.configs {
				db := manager.openConn(config)
				m[config.GetKey()] = db
			}
			manager.shardingDB = m
		}
	}

	return manager.shardingDB[key]
}

// lazy init
func (manager *TGormDataSourceManager) Master() *gorm.DB {

	if manager.masterDB == nil {
		l.Lock()
		defer l.Unlock()
		if manager.masterDB == nil {
			if manager.configs == nil {
				manager.configs = datasource.DataSourceConfigureInstance(server.DATASOURCE_CONFIGURE)
			}
			config := manager.configs[0]
			db := manager.openConn(config)
			manager.masterDB = db
		}
	}

	return manager.masterDB
}

func (manager *TGormDataSourceManager) Slave() *gorm.DB {

	if len(manager.configs) == 1 {
		return manager.Master()
	}

	if manager.slaveDB == nil {
		l.Lock()
		defer l.Unlock()
		if manager.slaveDB == nil {
			if manager.configs == nil {
				manager.configs = datasource.DataSourceConfigureInstance(server.DATASOURCE_CONFIGURE)
			}
			config := manager.configs[1]
			db := manager.openConn(config)
			manager.slaveDB = db
		}
	}

	return manager.slaveDB
}

func (manager *TGormDataSourceManager) openConn(config datasource.IDataSourceConfigure) *gorm.DB {
	db, err := gorm.Open(config.GetDriver(), config.GetURL())
	if err != nil {
		Log.WithFields(logrus.Fields{
			"username": config.GetUsername(),
			"password": config.GetPassword(),
			"url":      config.GetURL(),
			"driver":   config.GetDriver(),
		}).Error("DataSourceManager init error")
		if error := manager.tryToCreateDatabase(); error != nil {
			panic(error)
		}
	}

	db.DB().SetMaxIdleConns(config.GetIdlePoolSize())
	db.DB().SetMaxOpenConns(config.GetMaxPoolSize())
	db.DB().SetConnMaxLifetime(time.Duration(config.GetMaxLifeTime()) * time.Second)

	// 设置字符编码
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8")
	db.SingularTable(true)
	if config.GetSqlDebug() == 1 {
		db.LogMode(true)
		db.SetLogger(TSQLLogger{})
	}

	if config.AutoCreateTable() {
		// 如果有注册models, 则进行建表同步
		if len(manager.Models) > 0 {
			for _, m := range manager.Models {
				if !db.HasTable(m) {
					err := db.CreateTable(m).Error
					if err != nil {
						Log.Error(err)
					}
				}
			}
			db.AutoMigrate(manager.Models...)
		}

		//从config注册的model进行建表， 为了兼容老版本，将同时使用两种方式
		if len(config.GetModels()) > 0 {
			for _, m := range manager.Models {
				if !db.HasTable(m) {
					err := db.CreateTable(m).Error
					if err != nil {
						Log.Error(err)
					}
				}
			}
			db.AutoMigrate(config.GetModels()...)
		}
	}

	Log.WithField("db", db).Debug("create a new db connetion")
	return db
}

func (manager *TGormDataSourceManager) tryToCreateDatabase() error {
	Log.Info("Try to create a new db accrodding to url")
	for _, config := range manager.configs {
		driver := config.GetDriver()
		url := config.GetURL()

		if config.GetDriver() == "mysql" {
			var dbname string
			if strings.Contains(url, "?") {
				reg := regexp.MustCompile("/(.*)\\?")
				result := reg.FindStringSubmatch(url)
				if len(result) == 2 {
					dbname = string(result[1])
				}
			} else {
				dbname = strings.Split(url, "/")[1]
			}

			rootURL := strings.Replace(url, dbname, "", 1)
			db, error := gorm.Open(driver, rootURL)
			if error != nil {
				Log.WithFields(logrus.Fields{
					"dbname": dbname,
					"driver": driver,
					"url":    rootURL,
				}).Error("create database failed")
				return error
			}

			if error = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbname).Error; error != nil {
				Log.WithFields(logrus.Fields{
					"dbname": dbname,
					"driver": driver,
					"url":    rootURL,
				}).Error("create database failed")
				return error
			}
		}
	}
	return nil
}
