# shadow-framework
shaow-framework是一个简单的web框架，集成了github上比较流行的一些lib包和其他框架，主要是集成了gin，gorm，casbin，logrus，go-i18n等框架

gin 是一个web框架，请参考 https://github.com/gin-gonic/gin
gorm是一个orm框架，请参考 https://github.com/switfs/shadow-framework/orm/jinzhu/gorm
casbin是一个权限验证框架，请参考 https://github.com/casbin/casbin
logrus是一个log框架，请参考 https://github.com/sirupsen/logrus
go-i18n是一个国际化框架，请参考 https://github.com/nicksnyder/go-i18n


### best practice：
1. 不用使用git clone命令来下载任何源代码, 所有需要用到git clone命令的地方，全部用go get代替
2. 使用dep ensure -add添加依赖包，不要用git clone和 go get添加依赖包
3. 使用dep ensure -update更新依赖包，不要用git clone和 go get更新依赖包

quick start
===========

## 如何使用shadow-framework
1. shadow-framework使用dep来管理包依赖, dep安装参考https://github.com/golang/dep
2. 确保系统已经配置了GO_PATH，并且将GO_PATH/bin加入PATH
3. 在要引入shadow-framework的工程的根目录下执行dep init，此时会在当前目录生成一个vendor目录和Gopke.lock, Gopkg.toml文件

```
    dep init
```

4. 继续在当前目录执行go fetch命令，加入shadow-framework的依赖

    dep ensure -add github.com/switfs/shadow-framework

5. 此时在vendor目录下应该已经加入了所有依赖, 现在可以直接在代码中引用shadow

```go
import ( 
	"github.com/switfs/shadow-framework"
	"github.com/switfs/shadow-framework/security"
)

e := shadow.DefaultEngine()
e.GET("/", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "hello world",
    })
})
```
## 如何开发shadow-framework
1. 如果没有dep，请先安装dep
2. 设置GOPATH, 如果已经有GOPATH，可以直接使用，或者要新建GOPATH，请将最新的GOPATH置于其他GOPATH之前
3. 使用env GIT_TERMINAL_PROMPT=1 go get github.com/switfs/shadow-framework命令下载代码到GOPATH
4. 进入$GOPATH/src/github.com/switfs/shadow-framework目录，运行dep ensure, dep会自动添加所有的依赖到vendor目录
5. 完毕


shadow-framework interface
==========================
shadow-framework 提供了很多扩展点可以使使用者方便的进行功能扩展

## 数据库配置
默认数据库配置是用户提供一个配置文件，默认的路径为项目根目录下的./orm/config/datasource.json

```json
{
    "username": "root",
    "password": "111111",
    "url": "root:111111@tcp(127.0.0.1:3306)/casbin?charset=utf8&parseTime=True&loc=Local",
    "driver": "mysql"
}
```
如果想改变配置文件路径或使用其他格式的配置文件可以自己实现一个Configure 接口的实现，然后将此实现注册进我们的框架

```go
type Configure interface {
	Username() string
	Password() string
	Url() string
	Driver() string
}

```
```go
RegisterConfigure("Configure", yourConfigureConstructor)
```