package sessionmiddleware

import (
	"encoding/gob"

	"github.com/switfs/shadow-framework/logger"
	"github.com/switfs/shadow-framework/security"
)

var (
	Log *logger.Logger
)

func init() {
	Log = logger.InitLog()
	Log.Info("session middleware init")
	gob.Register(&security.TAnonymousAuthenticationToken{})
	gob.Register(&security.TRequestAuthenticationToken{})
	gob.Register(&security.TUsernamePasswordAuthenticationToken{})
	gob.Register(&security.TWebAuthenticationDetails{})
	gob.Register(&security.TUser{})
}
