package sutils

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	. "github.com/switfs/shadow-framework/extend/global"
	. "github.com/switfs/shadow-framework/logger"

	cmap "github.com/orcaman/concurrent-map"
)

const (
	STOP = "__P:"
)

var gortid uint64 = 10000
var grmutex sync.Mutex

func GetGoId() uint64 {
	grmutex.Lock()
	defer grmutex.Unlock()
	gortid = gortid + 1
	id := gortid
	return id
}

type GoroutineChannel struct {
	Gid       uint64
	Name      string
	Msg       chan string
	StartTime time.Time
}

type GoroutineChannelMap struct {
	// mutex      sync.Mutex
	// Grchannels map[string]*GoroutineChannel
	Grchannels cmap.ConcurrentMap
}

func (m *GoroutineChannelMap) unregister(name string) error {
	// m.mutex.Lock()
	// defer m.mutex.Unlock()
	if _, ok := m.Grchannels.Get(name); !ok {
		return fmt.Errorf("goroutine channel not find: %q", name)
	}
	// delete(m.grchannels, name)
	m.Grchannels.Remove(name)

	return nil
}

func (m *GoroutineChannelMap) register(name string, gid uint64) error {

	gchannel := &GoroutineChannel{
		// gid:  uint64(rand.Int63()),
		Gid:       gid,
		Name:      name,
		StartTime: time.Now(),
	}
	gchannel.Msg = make(chan string)
	// m.mutex.Lock()
	// defer m.mutex.Unlock()
	// if m.grchannels == nil {
	// 	// m.grchannels = make(map[string]*GoroutineChannel)
	// 	m.grchannels = cmap.New()
	// } else
	ASSERT(m.Grchannels != nil)

	if _, ok := m.Grchannels.Get(gchannel.Name); ok {
		return fmt.Errorf("goroutine channel already defined: %q", gchannel.Name)
	}
	// m.grchannels[gchannel.name] = gchannel
	m.Grchannels.Set(gchannel.Name, gchannel)
	return nil
}

type GoRoutineManager struct {
	GrchannelMap *GoroutineChannelMap
}

func NewGoRoutineManager() *GoRoutineManager {
	gm := &GoroutineChannelMap{
		Grchannels: cmap.New(),
	}
	return &GoRoutineManager{GrchannelMap: gm}
}

func (gm *GoRoutineManager) StopLoopGoroutine(name string) error {
	// stopChannel, ok := gm.grchannelMap.grchannels[name]
	obj, ok := gm.GrchannelMap.Grchannels.Get(name)
	if !ok {
		return fmt.Errorf("not found goroutine name :" + name)
	}
	stopChannel := obj.(*GoroutineChannel)

	// gm.grchannelMap.grchannels[name].msg <- STOP + strconv.Itoa(int(stopChannel.gid))
	stopChannel.Msg <- STOP + strconv.Itoa(int(stopChannel.Gid))

	return nil
}

func (gm *GoRoutineManager) NewLoopGoRoutine(name string, fc interface{}, args ...interface{}) {
	go func(this *GoRoutineManager, name string, fc interface{}, args ...interface{}) {

		defer CatchExceptionWithName(name)

		gid := GetGoId()
		name = fmt.Sprintf("%v-%v", name, gid)

		//register channel
		err := this.GrchannelMap.register(name, gid)
		if err != nil {
			Log.Errorf("register new loop goroutine[%v] failed, err = %v", name, err)
			return
		}

		obj, ok := this.GrchannelMap.Grchannels.Get(name)
		if !ok {
			panic("cannot found goroutine name :" + name)
		}
		stopChannel := obj.(*GoroutineChannel)

		for {
			select {
			case info := <-stopChannel.Msg:
				taskInfo := strings.Split(info, ":")
				signal, gid := taskInfo[0], taskInfo[1]
				if gid == strconv.Itoa(int(stopChannel.Gid)) {
					if signal == "__P" {
						fmt.Println("gid[" + gid + "] quit")
						this.GrchannelMap.unregister(name)
						return
					} else {
						fmt.Println("unknown signal")
					}
				}
			default:
				// fmt.Println("no signal")
			}

			if len(args) > 1 {
				fc.(func(...interface{}))(args...)
			} else if len(args) == 1 {
				fc.(func(interface{}))(args[0])
			} else {
				fc.(func())()
			}
		}
	}(gm, name, fc, args...)
}

func (gm *GoRoutineManager) NewGoRoutine(name string, fc interface{}, args ...interface{}) {
	go func(name string, fc interface{}, args ...interface{}) {

		defer CatchExceptionWithName(name)

		gid := GetGoId()
		name = fmt.Sprintf("%v-%v", name, gid)

		//register channel
		err := gm.GrchannelMap.register(name, gid)
		if err != nil {
			Log.Error("register new goroutine[%v] failed, err = %v", name, err)
			return
		}
		if len(args) > 1 {
			fc.(func(...interface{}))(args...)
		} else if len(args) == 1 {
			fc.(func(interface{}))(args[0])
		} else {
			fc.(func())()
		}
		gm.GrchannelMap.unregister(name)
	}(name, fc, args...)

}

func (gm *GoRoutineManager) PrintGoRoutineMap() {

	results := ""

	grmap := gm.GrchannelMap.Grchannels
	for item := range grmap.IterBuffered() {
		val := item.Val
		channel := val.(*GoroutineChannel)

		//Log.Info("goroutine channel = %v", channel.Name)
		results = results + fmt.Sprintf("%v\n", channel.Name)
	}

	Log.Info("==========size=%v\n%v\n=============", grmap.Count(), results)
}
