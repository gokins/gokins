
mkdir bin
wget -O bin/gokins-alpine http://down.gokins.cn/static/golang/linux64/gokins-alpine
chmod +x bin/gokins-alpine
docker build -t mgr9525/gokins:latest .
