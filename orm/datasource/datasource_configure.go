package datasource

type IDataSourceConfigure interface {
	GetKey() string
	GetUsername() string
	GetPassword() string
	GetURL() string
	GetDriver() string
	GetIdlePoolSize() int
	GetMaxPoolSize() int
	GetMaxLifeTime() int64
	GetSqlDebug() int8
	AutoCreateTable() bool
	GetModels() []interface{}
}
type FDataSourceConfigureFactory func() []IDataSourceConfigure

var DataSourceConfigureFactories = make(map[string]FDataSourceConfigureFactory)

func RegisterDataSourceConfigure(name string, factory FDataSourceConfigureFactory) {
	DataSourceConfigureFactories[name] = factory
}

func DataSourceConfigureInstance(name string) []IDataSourceConfigure {
	factory := DataSourceConfigureFactories[name]
	return factory()
}
