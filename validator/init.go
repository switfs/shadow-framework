package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/switfs/shadow-framework/logger"
)

var Log *logger.Logger

func init() {
	Log = logger.InitLog()
	Log.Infoln("Validator init")

}

func init() {
	binding.Validator = new(defaultValidator)
}
