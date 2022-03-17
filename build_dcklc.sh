
mkdir bin
wget -O bin/gokins-alpine http://down.gokins.cn/static/golang/linux64/gokins-alpine

docker build -t mgr9525/gokins:latest .
