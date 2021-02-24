package main

var TranslateHTML = `<div id='op'><button onclick="ToScroll()">滚动</button><button onclick="downloadF()">下载</button><input id="lang" value="请输入文件名"/></div>
<div id="content">$content$</div>
<script src="/public/api/apis.js"></script>
<script>
    window.ToScroll = function() {
        var oneTime = window.innerHeight / 2;
        var total = document.body.offsetHeight;
        var cur = 0;
        var id = setInterval(function () {
            if (cur < total) {
                cur += oneTime;
                window.scrollTo(0, cur);
            } else {
                clearInterval(id);
                window.scrollTo(0, 0);
                alert("完成");
            }
        }, 200)
    };

    window.downloadF = function() {
        var content = [],
            filename = lang.value + ".json";
        var div = document.getElementsByTagName('div');
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
</script>`
