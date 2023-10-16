package progconfig

import (
	"strconv"
	"strings"

	"github.com/switfs/shadow-framework/extend/sutils"
	"github.com/switfs/shadow-framework/orm/datasource"

	"github.com/switfs/shadow-framework/extend/base"
	. "github.com/switfs/shadow-framework/extend/global"

	. "github.com/switfs/shadow-framework/logger"

	"time"

	cmap "github.com/orcaman/concurrent-map"
)

type TModel struct {
	ID        int64 `gorm:"AUTO_INCREMENT;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TProgConfig struct {
	TModel
	ParamName  string `gorm:"size:64"`
	ParamValue string `gorm:"type:longtext"`
	Type       string `gorm:"size:32"`
	Disabled   bool   `gorm:""`
	Encrypted  bool   `gorm:""`
	Comment    string `gorm:"size:1024"`
}

const (
	CFG_DEF_PROG_CONFIG_RELOAD_INTERVAL = "DEF_PROG_CONFIG_RELOAD_INTERVAL"
	CFG_DEF_REDIS_ADDR                  = "DEF_REDIS_ADDR"
	CFG_DEF_REDIS_PASSWORD              = "DEF_REDIS_PASSWORD"
	CFG_DEF_REDIS_DB                    = "DEF_REDIS_DB"
	CFG_DEF_SESSION_DB                  = "DEF_SESSION_DB"
	CFG_DEF_CASBIN_CONFIG_FILE          = "DEF_CASBIN_CONFIG_FILE"
	CFG_DEF_NSQ_ADDR                    = "DEF_NSQ_ADDR"
	CFG_DEF_NSQ_HTTP_ADDR               = "DEF_NSQ_HTTP_ADDR"
	CFG_DEF_NSQ_ACK_TIMEOUT             = "DEF_NSQ_ACK_TIMEOUT"
	CFG_DEF_ENABLE_SESSION_CHECK        = "DEF_ENABLE_SESSION_CHECK"
	CFG_DEF_ENABLE_PRIVILEGE_CHECK      = "DEF_ENABLE_PRIVILEGE_CHECK"
	// CFG_DEF_METHOD_WITHOUT_SESSION_CHECK = "DEF_METHOD_WITHOUT_SESSION_CHECK"
	CFG_DEF_LOAD_LIBS = "DEF_LOAD_LIBS"
)

var defaultMap = map[string]string{
	CFG_DEF_PROG_CONFIG_RELOAD_INTERVAL: "600",
	CFG_DEF_REDIS_ADDR:                  "192.168.192.145:6379",
	CFG_DEF_REDIS_PASSWORD:              "",
	CFG_DEF_REDIS_DB:                    "0",
	CFG_DEF_SESSION_DB:                  "6",
	CFG_DEF_CASBIN_CONFIG_FILE:          "./config/casbin.conf",
	CFG_DEF_NSQ_ADDR:                    "192.168.192.145:4150",
	CFG_DEF_NSQ_HTTP_ADDR:               "http://192.168.192.145:4151",
	CFG_DEF_NSQ_ACK_TIMEOUT:             "120",
	CFG_DEF_ENABLE_SESSION_CHECK:        "1",
	CFG_DEF_ENABLE_PRIVILEGE_CHECK:      "1",
	// CFG_DEF_METHOD_WITHOUT_SESSION_CHECK: makeMethodsWithoutSessionCheck(),
	CFG_DEF_LOAD_LIBS: makeLoadLibs(),
}

type TProgConfigItem struct {
	value     interface{}
	disabled  bool
	encrypted bool
}

func makeLoadLibs() string {

	libs := []string{
		"./pluginso/public.so",
		"./pluginso/game.so",
		"./pluginso/apiv2.so",
		"./pluginso/zan.so",
	}

	return strings.Join(libs, ",")
}

// func makeMethodsWithoutSessionCheck() string {

// 	methods := []string{
// 		"/zan/seller/login",
// 		"/zan/customer/login",
// 		"/zan/customer/store/list",
// 		"/zan/customer/item/category/list",
// 		"/zan/customer/item/list",
// 		"/zan/customer/item/detail",
// 	}

// 	return strings.Join(methods, ",")

// }

func (item *TProgConfigItem) String() string {
	return item.value.(string)
}

func (item *TProgConfigItem) Int64() int64 {

	sval := item.value.(string)
	i, err := strconv.ParseInt(sval, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func (item *TProgConfigItem) Float64() float64 {

	sval := item.value.(string)

	f, err := strconv.ParseFloat(sval, 64)
	if err != nil {
		panic(err)
	}

	return f
}

func (item *TProgConfigItem) Decrypt() string {

	d, err := sutils.InternalDecryptStr(item.String())
	if err != nil {
		Log.Errorf("decrypt failed, err = %v", err)
		return ""
	}

	return d
}

func (item *TProgConfigItem) Disabled() bool {
	return item.disabled
}

type TProgConfigure struct {
	configMap cmap.ConcurrentMap
}

var progConfigure *TProgConfigure = &TProgConfigure{
	configMap: cmap.New(),
}

// var progConfigMap cmap.ConcurrentMap = cmap.New()

func ProgConfigureInit() {
	setConfig(defaultMap)

	loadProgConfig()

	manager := base.GetGoRoutineManager()
	manager.NewLoopGoRoutine("reloadProgConfig_loop", reloadProgConfig)
}

func ProgSetConfig(configmap map[string]string) {
	setConfig(configmap)
}

func setConfig(configmap map[string]string) {

	for k, v := range configmap {

		item := &TProgConfigItem{
			value:    v,
			disabled: false,
		}

		progConfigure.configMap.Set(k, item)
	}
}

func loadProgConfig() {
	var configs []TProgConfig

	db := datasource.DatasourceInstance().Master()
	err := db.Find(&configs).Error
	if err != nil {
		Log.Errorf("load prog config failed, err = %v", err)
		return
	}

	for _, c := range configs {

		item := &TProgConfigItem{
			value:     c.ParamValue,
			disabled:  c.Disabled,
			encrypted: c.Encrypted,
		}

		progConfigure.configMap.Set(c.ParamName, item)
	}
}

func reloadProgConfig() {

	loadProgConfig()

	intval := progConfigure.GetConfigItem(CFG_DEF_PROG_CONFIG_RELOAD_INTERVAL).Int64()

	time.Sleep(time.Duration(intval) * time.Second)
}

func ProgConfigureInstance() *TProgConfigure {
	ASSERT(progConfigure != nil)
	return progConfigure
}

func (p *TProgConfigure) GetConfigItem(name string) *TProgConfigItem {
	itemobj, ok := p.configMap.Get(name)
	ASSERT(ok)

	return itemobj.(*TProgConfigItem)
}

func (p *TProgConfigure) GetConfigItemDb(name string) *TProgConfigItem {
	config := TProgConfig{}

	db := datasource.DatasourceInstance().Master()
	err := db.Where("param_name=?", name).Find(&config).Error
	if err != nil {
		Log.Error("db error = %v", err)
		panic(err)
	}

	return &TProgConfigItem{value: config.ParamValue, disabled: config.Disabled}
}

func (p *TProgConfigure) SetConfigItemWithEncryption(name string, value string) error {

	encrypted, err := sutils.InternalEncryptStr(value)
	if err != nil {
		Log.Errorf("value=%v, err=%v", value, err)
	}

	return p.SetConfigItem(name, encrypted)
}

func (p *TProgConfigure) SetConfigItem(name string, value string) error {

	config := &TProgConfig{
		ParamName:  name,
		ParamValue: value,
	}

	db := datasource.DatasourceInstance().Master()
	err := db.Save(config).Error
	if err != nil {
		Log.Errorf("db error = %v", err)
		return err
	}

	p.configMap.Set(name, config)

	return nil
}

func GetConfigItem(name string) *TProgConfigItem {
	return ProgConfigureInstance().GetConfigItem(name)
}

func SetConfigItem(name string, value string) error {
	return ProgConfigureInstance().SetConfigItem(name, value)
}

func GetConfigItemDb(name string) *TProgConfigItem {
	return ProgConfigureInstance().GetConfigItemDb(name)
}

func SetConfigItemWithEncryption(name string, value string) error {
	return ProgConfigureInstance().SetConfigItemWithEncryption(name, value)
}
