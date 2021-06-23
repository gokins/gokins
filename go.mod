module github.com/gokins-main/gokins

go 1.15

require (
	github.com/alcortesm/tgz v0.0.0-20161220082320-9c5fe88206d7 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20210208195552-ff826a37aa15 // indirect
	github.com/gin-gonic/gin v1.7.2
	github.com/go-git/go-git/v5 v5.4.2 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gokins-main/core v0.0.0-20210622095351-4e4f8ce44617
	github.com/gokins-main/runner v0.0.0-00010101000000-000000000000
	github.com/kr/text v0.2.0 // indirect
	github.com/mgr9525/HyperByte-Transfer-Protocol v1.1.5
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/net v0.0.0-20210326060303-6b1517762897 // indirect
	golang.org/x/sys v0.0.0-20210502180810-71e4cd670f79 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.3.0
	xorm.io/xorm v1.1.0
)

replace (
	github.com/gokins-main/core => ../core
	github.com/gokins-main/runner => ../runner
)
