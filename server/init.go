package server

import (
	"github.com/switfs/shadow-framework/logger"
	"github.com/switfs/shadow-framework/orm/datasource"
)

var Log *logger.Logger

const (
	DATASOURCE_CONFIGURE = "DataSourceConfigure"
)

func init() {
	Log = logger.InitLog()
	Log.Infoln("ServerConfigure init")

	datasource.RegisterDataSourceConfigure(DATASOURCE_CONFIGURE, newDataSourceConfigure)
}
