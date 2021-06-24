package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/core"
	"github.com/gokins-main/core/common"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/migrates"
	"github.com/gokins-main/gokins/util"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
)

type installConfig struct {
	Server struct {
		Host     string `json:"host"` //外网访问地址
		RunLimit int    `json:"runLimit"`
		HbtpHost string `json:"hbtpHost"`
		Secret   string `json:"secret"`
		NoRun    bool   `json:"noRun"`
	} `json:"server"`
	Datasource struct {
		Driver string `json:"driver"`
		Host   string `json:"host"`
		Name   string `json:"name"`
		User   string `json:"user"`
		Pass   string `json:"pass"`
	} `json:"datasource"`
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
	if m.Datasource.Driver == "mysql" {
		if !common.RegHost2.MatchString(m.Datasource.Host) {
			c.String(500, "dbhost err:%s", m.Datasource.Host)
			return
		}
		if m.Datasource.Name == "" {
			c.String(500, "dbname err:%s", m.Datasource.Name)
			return
		}
		if strings.Contains(m.Datasource.Name, ":") || strings.Contains(m.Datasource.Pass, ":") {
			c.String(500, "(dbname & dbport) can't contains ':'")
			return
		}
	} else {
		m.Datasource.Driver = "sqlite"
	}

	dbwait := false
	dataul := ""
	var err error
	if m.Datasource.Driver == "mysql" {
		dbwait, dataul, err = migrates.InitMysqlMigrate(m.Datasource.Host, m.Datasource.Name, m.Datasource.User, m.Datasource.Pass)
	}
	if err != nil {
		if dbwait {
			c.String(512, "wait")
		} else {
			c.String(500, "init db err:%v", err)
		}
		return

	}
	if dataul == "" {
		c.String(513, "datasource info err")
		return
	}

	comm.Cfg.Server.Host = m.Server.Host
	comm.Cfg.Server.Shells = []string{"shell@sh", "shell@bash"}
	if runtime.GOOS == "windows" {
		comm.Cfg.Server.Shells = []string{"shell@cmd", "shell@powershell"}
	}
	if m.Server.NoRun {
		comm.Cfg.Server.Shells = nil
	}
	comm.Cfg.Server.HbtpHost = m.Server.HbtpHost
	comm.Cfg.Server.Secret = m.Server.Secret
	comm.Cfg.Datasource.Driver = m.Datasource.Driver
	comm.Cfg.Datasource.Url = dataul
	err = initConfig()
	if err != nil {
		c.String(500, "init config err:%v", err)
		return
	}
	comm.Installed = true
	c.String(200, "ok")
}

func initConfig() error {
	bts, err := yaml.Marshal(&comm.Cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(comm.WorkPath, "app.yml"), bts, 0644)
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
        
        .content div {
            margin-bottom: 10px;
        }
        
        #msgDiv {
            color: red;
        }
    </style>
</head>

<body>
    <div class="content">
        <div>访问地址:<input id="hostTxt" style="width: 500px;" /></div>
        <div>插件服务:
            <select id="plugServ">
            <option value="">不启用</option>
            <option value="1">内网</option>
            <option value="2">外网</option>
          </select>
            <input id="plugPort" style="width: 100px;" placeholder="端口" />
            <input id="plugSecret" style="width: 300px;" placeholder="密码" />
        </div>
        <div>数据库:
            <select id="dbDriver">
              <option value="sqlite">sqlite</option>
          <option value="mysql">mysql</option>
        </select>

        </div>
        <div id="mysqlDiv">
            <div>dbhost:<input id="dbhostTxt" style="width: 500px;" value="localhost:3306" /></div>
            <div>dbname:<input id="dbnameTxt" style="width: 500px;" value="gokins" /></div>
            <div>dbuser:<input id="dbuserTxt" style="width: 500px;" value="root" /></div>
            <div>dbpass:<input id="dbpassTxt" style="width: 500px;" /></div>
        </div>
        <div id="msgDiv"></div>
        <div style="margin-top: 20px;"><button id="subBtn" type="button" onclick="onInstal()">立即安装</button></div>
    </div>
    <script src="http://180.76.38.108:7039/test/axios.min.js"></script>
    <script src="http://180.76.38.108:7039/test/jquery-2.1.0.js"></script>
    <script>
        var msgDiv = $('#msgDiv');
        var subBtn = $('#subBtn');
        var service = axios.create({
            baseURL: "/api", // api base_url
            // baseURL: 'http://n.1ydt.com:8072', // api base_url
            //timeout: 5000, // 请求超时时间
            withCredentials: true
        });

        var regul = /^(https?:)\/\/([\w\.]+)(:\d+)?/;
        var reghost = /^([\w\.]+)(:\d+)?$/;

        function plugChange() {
            switch ($('#plugServ').val()) {
                case '1':
                    $('#plugSecret').val('');
                    $('#plugPort').removeAttr('disabled');
                    $('#plugSecret').prop('disabled', 'disabled');
                    break
                case '2':
                    $('#plugPort').removeAttr('disabled');
                    $('#plugSecret').removeAttr('disabled');
                    break
                default:
                    $('#plugPort').val('');
                    $('#plugSecret').val('');
                    $('#plugPort').prop('disabled', 'disabled');
                    $('#plugSecret').prop('disabled', 'disabled');
                    break
            }
        }
        plugChange()
        $('#plugServ').on('change', plugChange);

        function dbChange() {
            switch ($('#dbDriver').val()) {
                case 'mysql':
                    $('#mysqlDiv').show();
                    break
                default:
                    $('#mysqlDiv').hide();
                    break
            }
        }
        dbChange()
        $('#dbDriver').on('change', dbChange);

        function onInstal() {
            try {
                var csjs = {
                    "server": {
                        "host": $('#hostTxt').val()
                    },
                    "datasource": {
                        "driver": ''
                    }
                };
                if (!regul.test(csjs.server.host)) {
                    msgDiv.text('参数错误:host格式错误');
                    return
                }
                switch ($('#plugServ').val()) {
                    case '1':
                        var hbtpPort = $('#plugPort').val();
                        if (!/^\d+$/.test(hbtpPort)) {
                            msgDiv.text('参数错误:port格式错误');
                            return
                        }
                        csjs.server.hbtpHost = '127.0.0.1:' + hbtpPort;
                        break
                    case '2':
                        var hbtpPort = $('#plugPort').val();
                        if (!/^\d+$/.test(hbtpPort)) {
                            msgDiv.text('参数错误:port格式错误');
                            return
                        }
                        csjs.server.hbtpHost = ':' + hbtpPort;
                        csjs.server.secret = $('#plugSecret').val();
                        break
                }
                switch ($('#dbDriver').val()) {
                    case 'mysql':
                        var dburl = '';
                        var dbhost = $('#dbhostTxt').val();
                        var dbname = $('#dbnameTxt').val();
                        var dbuser = $('#dbuserTxt').val();
                        var dbpass = $('#dbpassTxt').val();
                        if (!reghost.test(dbhost)) {
                            msgDiv.text('参数错误:dbhost格式错误');
                            return
                        }
                        if (dbname == '') {
                            msgDiv.text('参数错误:dbname必填');
                            return
                        }
                        if (dbuser == '') {
                            msgDiv.text('参数错误:dbuser必填');
                            return
                        }
                        csjs.datasource.driver = 'mysql';
                        csjs.datasource.host = dbhost;
                        csjs.datasource.name = dbname;
                        csjs.datasource.user = dbuser;
                        csjs.datasource.pass = dbpass;
                        break
                    default:
                        csjs.datasource.driver = 'sqlite';
                        break
                }

                console.log('start install', csjs);
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
            $('#hostTxt').val(mts[0]);
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
