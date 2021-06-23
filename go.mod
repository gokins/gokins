module github.com/gokins-main/gokins

go 1.15

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20210208195552-ff826a37aa15 // indirect
	github.com/gin-gonic/gin v1.7.2
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gokins-main/core v0.0.0-20210622095351-4e4f8ce44617
	github.com/mgr9525/HyperByte-Transfer-Protocol v1.1.5
	github.com/sirupsen/logrus v1.8.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v2 v2.2.8
	xorm.io/xorm v1.1.0
)

replace github.com/gokins-main/core => ../core
