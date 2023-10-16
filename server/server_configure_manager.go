package server

import (
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/switfs/shadow-framework/orm/datasource"
)

var (
	DefaultConfigurePath string = "config/server.json"
	serverConfigure      IServerConfigure
)

type TServerConfigure struct {
	ServerName string
	Platform   string
	Node       int
	DataSource []TDataSourceConfig
}

type IServerConfigure interface {
	GetDataSource() []datasource.IDataSourceConfigure
}

func (c *TServerConfigure) GetDataSource() []datasource.IDataSourceConfigure {
	ret := make([]datasource.IDataSourceConfigure, len(c.DataSource))
	for index, config := range c.DataSource {
		config.AutoCreate = true
		ret[index] = &config
	}
	return ret
}

func ServerConfigureInstance() IServerConfigure {
	return newServerConfigure()
}

func newServerConfigure() IServerConfigure {
	if serverConfigure == nil {
		config := &TServerConfigure{}
		LoadWithFile(config, DefaultConfigurePath)
		Log.WithFields(logrus.Fields{
			"ServerName": config.ServerName,
			"Platform":   config.Platform,
			"Node":       config.Node,
		}).Debug("TServerConfigure")
		serverConfigure = config
	}

	return serverConfigure
}

func newDataSourceConfigure() []datasource.IDataSourceConfigure {
	return newServerConfigure().GetDataSource()
}

func LoadWithFile(configure interface{}, path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		Log.WithField("file", path).Error("server configure init failed, file doesn't exist")
		Log.Panic(err)
	}
	Log.Info("Server configure", string(data))

	datajson := []byte(data)
	err = json.Unmarshal(datajson, configure)
	if err != nil {
		Log.Panic(err)
	}
	if config, ok := configure.(IServerConfigure); ok {
		serverConfigure = config
	}
}

type TDataSourceConfig struct {
	Key          string
	Username     string
	Password     string
	URL          string
	Driver       string
	IdlePoolSize int
	MaxPoolSize  int
	MaxLifeTime  int64
	SqlDebug     int8
	AutoCreate   bool
	Models       []interface{}
}

func (configure *TDataSourceConfig) GetKey() string {
	return configure.Key
}

func (configure *TDataSourceConfig) GetUsername() string {
	return configure.Username
}

func (configure *TDataSourceConfig) GetPassword() string {
	return configure.Password
}

func (configure *TDataSourceConfig) GetURL() string {
	return configure.URL
}

func (configure *TDataSourceConfig) GetDriver() string {
	return configure.Driver
}

func (configure *TDataSourceConfig) GetIdlePoolSize() int {
	return configure.IdlePoolSize
}

func (configure *TDataSourceConfig) GetMaxPoolSize() int {
	return configure.MaxPoolSize
}

func (configure *TDataSourceConfig) GetMaxLifeTime() int64 {
	return configure.MaxLifeTime
}

func (configure *TDataSourceConfig) GetSqlDebug() int8 {
	return configure.SqlDebug
}

func (configure *TDataSourceConfig) AutoCreateTable() bool {
	return configure.AutoCreate
}

func (configure *TDataSourceConfig) GetModels() []interface{} {
	return configure.Models
}

func (configure *TDataSourceConfig) RegisterModels(models ...interface{}) {
	configure.Models = append(configure.Models, models...)
}
