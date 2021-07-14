package route

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gokins-main/gokins/util/httpex"

	"github.com/gin-gonic/gin"
	"github.com/gokins-main/core"
	"github.com/gokins-main/core/common"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/migrates"
	"github.com/gokins-main/gokins/util"
	"gopkg.in/yaml.v3"
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
	g.POST("/check", cs.check)
	g.POST("/", util.GinReqParseJson(cs.install))
}
func (InstallController) check(c *gin.Context) {
	c.String(200, "hello gokins!")
}
func checkUrl(host string) bool {
	req, err := http.NewRequest("POST", host+"/api/install/check", nil)
	if err != nil {
		return false
	}
	cli := http.Client{}
	cli.Timeout = time.Second * 5
	res, err := cli.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return res.StatusCode == 200
}
func (InstallController) install(c *gin.Context, m *installConfig) {
	if strings.HasSuffix(m.Server.Host, "/") {
		ln := len(m.Server.Host)
		m.Server.Host = m.Server.Host[:ln-2]
	}
	if !common.RegUrl.MatchString(m.Server.Host) {
		c.String(500, "host err:%s", m.Server.Host)
		return
	}
	if !checkUrl(m.Server.Host) {
		c.String(511, "can't connect:%s", m.Server.Host)
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

	dataul := ""
	var err error
	if m.Datasource.Driver == "mysql" {
		_, dataul, err = migrates.InitMysqlMigrate(m.Datasource.Host, m.Datasource.Name, m.Datasource.User, m.Datasource.Pass)
	} else {
		dataul, err = migrates.InitSqliteMigrate()
	}
	if err != nil {
		c.String(512, "%v", err)
		return

	}
	if dataul == "" {
		c.String(513, "datasource info err")
		return
	}

	comm.Cfg.Server.Host = m.Server.Host
	comm.Cfg.Server.LoginKey = utils.RandomString(32)
	comm.Cfg.Server.DownToken = utils.RandomString(32)
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
		httpex.ResMsgUrl(c, "重复操作,跳转中...", "/")
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
				left: 50%;
				left: 50%;
				top: 10%;
				width: 800px;
				margin-left: -400px;
				background: #eee;
				padding: 10px;
			}
	
			.content .layui-card div {
				margin-bottom: 10px;
			}
	
			.headers {
				background: #fff;
				padding-bottom: 10px;
				line-height: 70px;
				display: flex;
			}
	
			.headers h1 {
				color: #5360f7;
			}
	
			.login-dog {
				width: 100px;
				height: 60px;
				margin: 10px 20px 0 0;
			}
	
			#msgDiv {
				color: red;
				text-align: center;
			}
		</style>
		<!-- 引入 layui.css -->
		<link rel="stylesheet" href="gokinsui/plugs/layui.css" />
		<!-- 引入 layui.js -->
	</head>
	
	<body>
		<div class="content">
			<div class="headers" style="margin: 0;">
				<div style="flex: 1;"></div>
				<img class="login-dog" src="gokinsui/imgs/logo.png" />
				<div style="padding-top: 5px;">
					<h1>安装Gokins</h1>
				</div>
				<div style="flex: 1;"></div>
			</div>
	
			<div class="layui-card">
				<div class="layui-card-body">
					<form class="layui-form" action="javascript:onInstal()">
						<div class="layui-form-item">
							<label class="layui-form-label">访问地址</label>
							<div class="layui-input-inline" style="width:300px">
								<input type="text" id="hostTxt" name="txt1" lay-verify="required" placeholder="请输入标题" autocomplete="off"
									class="layui-input">
							</div>
							<div class="layui-form-mid layui-word-aux">需要 webhook,ssh,制品下载 建议使用<span style="color:red">外网</span></div>
						</div>
						<div class="layui-form-item">
							<label class="layui-form-label">插件服务</label>
							<div class="layui-input-inline">
								<select id="plugServ" lay-filter="plugs">
									<option value="">不启用</option>
									<option value="1">内网</option>
									<option value="2">外网</option>
								</select>
							</div>
							<div class="layui-input-inline" style="width:100px">
								<input type="text" id="plugPort" name="txt2" placeholder="服务端口" autocomplete="off" class="layui-input">
							</div>
							<div class="layui-input-inline" style="width:200px">
								<input type="text" id="plugSecret" name="txt3" placeholder="服务Secret" autocomplete="off"
									class="layui-input">
							</div>
							<!-- <div class="layui-form-mid layui-word-aux">差</div> -->
						</div>
						<div class="layui-form-item">
							<label class="layui-form-label">数据库</label>
							<div class="layui-input-inline">
								<select id="dbDriver" disabled="disabled">
									<option value="sqlite">sqlite</option>
									<option value="mysql" selected>mysql</option>
								</select>
							</div>
						</div>
						<div class="layui-form-item">
							<label class="layui-form-label">数据库地址</label>
							<div class="layui-input-inline" style="width:300px">
								<input type="text" id="dbhostTxt" name="txt4" lay-verify="required" value="localhost:3306"
									autocomplete="off" class="layui-input">
							</div>
							<div class="layui-form-mid layui-word-aux">Mysql链接地址</div>
						</div>
						<div class="layui-form-item">
							<label class="layui-form-label">数据库名称</label>
							<div class="layui-input-inline" style="width:200px">
								<input type="text" id="dbnameTxt" name="txt5" lay-verify="required" value="gokins" autocomplete="off"
									class="layui-input">
							</div>
						</div>
						<div class="layui-form-item">
							<label class="layui-form-label">数据库用户</label>
							<div class="layui-input-inline" style="width:200px">
								<input type="text" id="dbuserTxt" name="txt5" required lay-verify="required" value="root"
									autocomplete="off" class="layui-input">
							</div>
						</div>
						<div class="layui-form-item">
							<label class="layui-form-label">数据库密码</label>
							<div class="layui-input-inline" style="width:200px">
								<input type="text" id="dbpassTxt" name="txt6" value="" autocomplete="off" class="layui-input">
							</div>
						</div>
						<div class="layui-form-item">
							<div style="text-align: center;">
								<button class="layui-btn layui-btn-normal" lay-submit lay-filter="formd" id="subBtn">立即安装</button>
							</div>
						</div>
					</form>
					<div id="msgDiv"></div>
				</div>
			</div>
		</div>
		<script src="gokinsui/plugs/axios.js"></script>
		<script src="gokinsui/plugs/jquery.js"></script>
		<script src="gokinsui/plugs/layui.js"></script>
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
						if($('#plugPort').val()=='')
							$('#plugPort').val('8031');
						$('#plugSecret').val('');
						$('#plugPort').removeAttr('disabled');
						$('#plugSecret').prop('disabled', 'disabled');
						break
					case '2':
						if($('#plugPort').val()=='')
							$('#plugPort').val('8031');
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
	
			/* function dbChange() {
				switch ($('#dbDriver').val()) {
					case 'mysql':
						$('#mysqlDiv').show();
						break
					default:
						$('#mysqlDiv').hide();
						break
				}
			}
			dbChange() */
	
	
			var lyform = null;
			layui.use('form', function () {
				lyform = layui.form;
				//各种基于事件的操作，下面会有进一步介绍
	
				lyform.on('select(plugs)', function (data) {
					plugChange()
				})
			});
	
	
	
	
	
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
						layer.msg('访问地址格式错误', { icon: 2 });
						return
					}
					switch ($('#plugServ').val()) {
						case '1':
							var hbtpPort = $('#plugPort').val();
							if (!/^\d+$/.test(hbtpPort)) {
								layer.msg('插件服务端口格式错误', { icon: 2 });
								return
							}
							csjs.server.hbtpHost = '127.0.0.1:' + hbtpPort;
							break
						case '2':
							var hbtpPort = $('#plugPort').val();
							if (!/^\d+$/.test(hbtpPort)) {
								layer.msg('插件服务端口格式错误', { icon: 2 });
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
								layer.msg('数据库地址格式错误', { icon: 2 });
								return
							}
							if (dbname == '') {
								layer.msg('数据库名称必填', { icon: 2 });
								return
							}
							if (dbuser == '') {
								layer.msg('数据库用户必填', { icon: 2 });
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
					subBtn.addClass('layui-btn-disabled');
					service.post('/install', csjs).then(function (res) {
						// subBtn.removeAttr('disabled');
						subBtn.value = '跳转中';
						msgDiv.text('安装成功:' + res.data);
						layer.msg('安装成功:' + res.data, { icon: 1 });
						setTimeout(function () {
							window.location = 'gokins';
						}, 1000)
					}).catch(function (err) {
						subBtn.removeAttr('disabled');
						subBtn.removeClass('layui-btn-disabled');
						console.log('install err:', err);
						msgDiv.text('安装失败:' + (err.response ? err.response.data || '服务器错误' : '网络错误'));
	
						const stat = err.response ? err.response.status : 0;
						switch (stat) {
							case 404:
								window.location = 'gokins';
								break
							case 511:
								layer.msg('无法连接访问地址,请重试!', { icon: 2 });
								break
							case 512:
								layer.msg('无法初始化数据库,请重试!', { icon: 2 });
								break
						}
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
