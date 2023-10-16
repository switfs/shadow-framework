package i18nmiddleware

import (
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/switfs/shadow-framework/logger"
)

var (
	Log *logger.Logger
)

func init() {
	Log = logger.InitLog()
	Log.Info("I18nResolver init")
	i18n.MustLoadTranslationFile("config/i18n/en.all.json")
	i18n.MustLoadTranslationFile("config/i18n/zh.all.json")
}
