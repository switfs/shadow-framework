package datasource

import (
	"github.com/switfs/shadow-framework/orm/jinzhu/gorm"
)

type IDatasourceManager interface {
	Datasource() *gorm.DB
	Master() *gorm.DB
	Slave() *gorm.DB
	Configs(configs []IDataSourceConfigure)
	RegisterModels(models ...interface{})
}
type FDatasourceManagerFactory func() IDatasourceManager

var datasourceManagerFactories = make(map[string]FDatasourceManagerFactory)

func RegisterDatasourceManager(name string, factory FDatasourceManagerFactory) {
	datasourceManagerFactories[name] = factory
}

func DatasourceManagerInstance(name string) IDatasourceManager {
	factory := datasourceManagerFactories[name]
	return factory()
}

func DatasourceInstance() IDatasourceManager {
	factory := datasourceManagerFactories[DATASOURCE_MANAGER]
	return factory()
}
