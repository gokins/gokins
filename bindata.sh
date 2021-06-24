
# go get -u github.com/jteeuwen/go-bindata/...

go-bindata -o comm/migrate.go -pkg=comm -prefix migrates migrates/mysql/ migrates/sqlite/