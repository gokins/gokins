module github.com/gokins/gokins

go 1.15

require (
	github.com/boltdb/bolt v1.3.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.7.2
	github.com/go-git/go-git/v5 v5.4.2
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gokins/core v0.0.0-20221114063511-e32796e39c37
	github.com/gokins/runner v0.0.0-20230926040038-6d063511474e
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/mgr9525/HyperByte-Transfer-Protocol v1.1.6-0.20221010061341-5cadd93a6fab
	github.com/sirupsen/logrus v1.8.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
	xorm.io/builder v0.3.8
	xorm.io/xorm v1.1.0
)

// replace github.com/gokins/runner => ../runner_dev
