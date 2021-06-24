package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/core"
	"github.com/gokins-main/core/common"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/util"
	"io/ioutil"
	"strings"
)

type installConfig struct {
	Server struct {
		Host     string `json:"host"` //外网访问地址
		RunLimit int    `json:"runLimit"`
		HbtpHost string `json:"hbtpHost"`
		Secret   string `json:"secret"`
	} `json:"server"`
	Database struct {
		Driver string `json:"driver"`
		Host   string `json:"host"`
		Name   string `json:"name"`
		Pass   string `json:"pass"`
	} `json:"database"`
}
type InstallController struct{}

func (InstallController) GetPath() string {
	return "/api/install"
}
func (cs *InstallController) auth(c *gin.Context) {
	if comm.Installed {
		c.String(404, "Not Found")
		c.Abort()
		return
	}
}
func (cs *InstallController) Routes(g gin.IRoutes) {
	g.Use(cs.auth)
	g.POST("/", util.GinReqParseJson(cs.install))
}
func (InstallController) install(c *gin.Context, m *installConfig) {
	if !common.RegUrl.MatchString(m.Server.Host) {
		c.String(500, "host err:%s", m.Server.Host)
		return
	}
	if m.Server.HbtpHost != "" && !common.RegHost1.MatchString(m.Server.HbtpHost) {
		c.String(500, "hbtp host err:%s", m.Server.HbtpHost)
		return
	}
	if m.Database.Driver == "mysql" {
		if !common.RegHost2.MatchString(m.Database.Host) {
			c.String(500, "dbhost err:%s", m.Database.Host)
			return
		}
		if m.Database.Name == "" {
			c.String(500, "dbname err:%s", m.Database.Name)
			return
		}
		if strings.Contains(m.Database.Name, ":") || strings.Contains(m.Database.Pass, ":") {
			c.String(500, "(dbname & dbport) can't contains ':'")
			return
		}
	} else {
		m.Database.Driver = "sqlite"
	}

}

func Install(c *gin.Context) {
	if comm.Installed {
		c.String(404, "Not Found")
		return
	}
	bts := []byte(`
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>安装</title>
    <style>
        .content {
            position: absolute;
            left: 0;
            right: 0;
            top: 30%;
            bottom: 0;
            text-align: center;
            font-size: 16px;
            color: #443ad6;
        }
        
        #msgDiv {
            color: red;
        }
    </style>
</head>

<body>
    <div class="content">
        <div>host:<input id="hostTxt" style="width: 500px;" /></div>
        <div>port:<input id="portTxt" style="width: 500px;" /></div>
        <div>dbhost:<input id="dbhostTxt" style="width: 500px;" value="180.76.38.108" /></div>
        <div>dbport:<input id="dbportTxt" style="width: 500px;" value="3306" /></div>
        <div>dbpass:<input id="dbpassTxt" style="width: 500px;" /></div>
        <div>是否启用内置runner:<input type="checkbox" id="isMainRun" checked="checked" /></div>
        <div id="msgDiv"></div>
        <div style="margin-top: 20px;"><button id="subBtn" type="button" onclick="onInstal()">立即安装</button></div>
    </div>
    <script src="http://180.76.38.108:7039/test/axios.min.js"></script>
    <script src="http://180.76.38.108:7039/test/jquery-2.1.0.js"></script>
    <script>
        var msgDiv = $('#msgDiv');
        var hostTxt = $('#hostTxt');
        var portTxt = $('#portTxt');
        var dbhostTxt = $('#dbhostTxt');
        var dbportTxt = $('#dbportTxt');
        var dbpassTxt = $('#dbpassTxt');
        var isMainRun = $('#isMainRun');
        var subBtn = $('#subBtn');
        var service = axios.create({
            baseURL: "/api", // api base_url
            // baseURL: 'http://n.1ydt.com:8072', // api base_url
            //timeout: 5000, // 请求超时时间
            withCredentials: true
        });

        var regul = /^(https?:)\/\/([\w\.]+)(:\d+)?/;

        function onInstal() {
            try {
                var csjs = {
                    "server": {
                        "host": hostTxt.val(),
                        "port": portTxt.val(),
                        "login-key": "",
                        limit: {
                            "build": 5,
                        }
                    },
                    "runner": {
                        "port": "",
                        "secret": "",
                        "main-runner": isMainRun.prop('checked')
                    },
                    "datasource": {
                        "host": dbhostTxt.val(),
                        "port": dbportTxt.val(),
                        "username": "root",
                        "password": dbpassTxt.val(),
                        "database": "gitee-go-dev"
                    },
                };
                if (!regul.test(csjs.server.host)) {
                    msgDiv.text('参数错误:host格式错误');
                    return
                }
                if (!/^\d+$/.test(csjs.server.port)) {
                    msgDiv.text('参数错误:port格式错误');
                    return
                }
                if (csjs.datasource.host == '' || csjs.datasource.password == '') {
                    msgDiv.text('参数错误:数据库密码未填写');
                    return
                }
                if (!/^\d+$/.test(csjs.datasource.port)) {
                    msgDiv.text('参数错误:dbport格式错误');
                    return
                }
                subBtn.attr('disabled', 'disabled');
                service.post('/install', csjs).then(function(res) {
                    subBtn.removeAttr('disabled');
                    console.log('install ok:', res.data);
                    msgDiv.text('安装成功:' + res.data);
                }).catch(function(err) {
                    subBtn.removeAttr('disabled');
                    console.log('install err:', err);
                    msgDiv.text('安装失败:' + err.response ? err.response.data || '服务器错误' : '网络错误');
                });
            } catch (e) {
                msgDiv.text('安装失败,json错误:' + e);
            }
        }

        var hrefs = window.location.href;
        if (regul.test(hrefs)) {
            var mts = hrefs.match(regul);
            console.log('match', mts);
            hostTxt.val(mts[0]);
            if (mts.length >= 3)
                portTxt.val(mts[3].replace(':', ''));
        }
    </script>
</body>

</html>
	`)
	if core.Debug {
		bs, err := ioutil.ReadFile("install.html")
		if err == nil {
			bts = bs
		}
	}
	c.Data(200, "text/html", bts)
}
