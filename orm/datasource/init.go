package datasource

import (
	"github.com/switfs/shadow-framework/logger"
)

var Log *logger.Logger

const (
	DATASOURCE_MANAGER = "DataSourceManager"
)

func init() {
	Log = logger.InitLog()
}
