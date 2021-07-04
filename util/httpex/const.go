package httpex

import (
	"github.com/gin-gonic/gin"
	"strings"
)

var HTMLMsgUrl = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>跳转提示</title>
    <style>
        #content {
            position: absolute;
            left: 0;
            right: 0;
            top: 45%;
            bottom: 0;
            text-align: center;
            font-size: 16px;
            color: #443ad6;
        }
    </style>
</head>

<body>
    <div id="content">{{msg}}</div>
    <script>
        var contDiv = document.getElementById("content");
        var msg = contDiv.innerHTML;
        var url = "{{url}}";
        if (msg == null || msg.length <= 0)
            contDiv.innerHTML = '跳转中...';
        setTimeout(function() {
            if (url != null && url.length > 0) {
                window.location.replace(url);
            } else
                window.location = '/';
        }, (msg != null && msg.length > 0) ? 1000 : 100);
    </script>
</body>
</html>
		`

func ResMsgUrl(c *gin.Context, msg string, url ...string) {
	hls := strings.ReplaceAll(HTMLMsgUrl, "{{msg}}", msg)
	if len(url) > 0 {
		hls = strings.ReplaceAll(hls, "{{url}}", url[0])
	} else {
		hls = strings.ReplaceAll(hls, "{{url}}", "")
	}
	c.Data(302, "text/html", []byte(hls))
}
