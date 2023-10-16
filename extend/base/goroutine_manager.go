package base

import (
	. "github.com/switfs/shadow-framework/extend/global"
	"github.com/switfs/shadow-framework/extend/sutils"
)

var goRoutineManager *sutils.GoRoutineManager = sutils.NewGoRoutineManager()

func GetGoRoutineManager() *sutils.GoRoutineManager {
	ASSERT(goRoutineManager != nil)
	return goRoutineManager
}
