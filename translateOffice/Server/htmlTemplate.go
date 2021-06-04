package Server

var TranslateHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
<div id='op'>
<button onclick="nextPage()">下一页</button>
<button onclick="ToScroll()">滚动</button>
<label>滚动速度(越右越慢)</label>
<input id="speed" style="margin:0 5px;" value="10" type="range" max="100" min="1" step="1"/>
<button onclick="downloadF()">下载</button>
<input id="lang" style="margin-left:5px;" value="请输入文件名"/></div>
<div id="content">
    $content$
</div>
<script>
	window.nextPage = function() {
		window.location = "/" + (+location.href.split('/')[3] + 1);
	}
	window.onload = function() {
		document.getElementById('lang').value = +location.href.split('/')[3];
		document.getElementById('speed').value = + (localStorage.getItem('speed') || 10);
	}
    window.ToScroll = function() {
        var total = document.body.offsetHeight;
        var cur = 0;
        var id = setInterval(function () {
			localStorage.setItem('speed', +document.getElementById('speed').value);
        	var oneTime = window.innerHeight / +document.getElementById('speed').value;
            if (cur < total) {
                cur += oneTime;
                window.scrollTo(0, cur);
            } else {
                clearInterval(id);
                window.scrollTo(0, 0);
                alert("完成");
            }
        }, 200);
    };

    window.downloadF = function() {
        var content = [],
            filename = document.getElementById('lang').value + ".json";
        var div = document.getElementById('content').getElementsByTagName('div');
        for (var i = 0; i < div.length; i++) {
            content.push({
                id: div[i].id,
                content: div[i].innerText
            });
        }
        content = JSON.stringify(content,'','\t');
        // 创建a标签
        let linkNode = document.createElement('a');
        linkNode.download = filename;
        linkNode.style.display = 'none';
        // 利用Blob对象将字符内容转变成二进制数据
        let blob = new Blob([content]);
        linkNode.href = URL.createObjectURL(blob);
        // 点击
        document.body.appendChild(linkNode);
        linkNode.click();
        // 移除
        document.body.removeChild(linkNode);
    };
</script>
</body>
</html>`
