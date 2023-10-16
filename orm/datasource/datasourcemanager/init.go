package datasourcemanager

import (
	"sync"

	"github.com/switfs/shadow-framework/logger"
	"github.com/switfs/shadow-framework/orm/datasource"
)

var (
	Log *logger.Logger
	l   sync.Mutex
)

func init() {
	Log = logger.InitLog()
	Log.Infoln("GormDataSourceManager init")
	datasource.RegisterDatasourceManager(datasource.DATASOURCE_MANAGER, newGormDataSourceManager)
	datasource.RegisterShardingDatasourceManager(datasource.DATASOURCE_MANAGER, newGormShardingDatasourceManager)
}
