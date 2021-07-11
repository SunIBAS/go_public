package main

import (
	"fmt"
	"net/http"
	"path"
	"public.sunibas.cn/go_public/utils/Console"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"strconv"
	"strings"
	"time"
)

var (
	basePath = "D:\\easyHelper"
	port     = ":10881"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		if len(url) >= len("/_video") && url[0:len("/_video")] == "/_video" {
			playVideo(w, r)
		} else if len(url) > 4 {
			pInd := strings.LastIndex(url, ".")
			if url[pInd:] == ".mp4" {
				http.ServeFile(w, r, path.Join(basePath, url))
			}
		} else {
			returnFileList(w)
		}
	})
	srv := &http.Server{
		Addr:           port,
		Handler:        nil,
		ReadTimeout:    time.Duration(5) * time.Minute,
		WriteTimeout:   time.Duration(5) * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	srv.ListenAndServe()
	Console.RunBatFile("http://localhost" + port)
}

//获取URL的GET参数
func GetUrlArg(r *http.Request, name string) string {
	var arg string
	values := r.URL.Query()
	arg = values.Get(name)
	return arg
}
func fzf(w http.ResponseWriter) {
	content := "404"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	fmt.Fprint(w, content)
}
func returnFileList(w http.ResponseWriter) {
	fn := DirAndFile.GetSubDirOrFile(basePath)
	content := ""
	for _, node := range fn {
		pInd := strings.LastIndex(node.Name, ".")
		if pInd > 0 {
			if node.Name[pInd:] == ".mp4" {
				content += `<div><a href="/_video?video=` + node.Name + `">` + node.Name + `</a></div>`
				continue
			}
		}
		content += `<div><a href="` + node.Name + `">` + node.Name + `</a></div>`
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	fmt.Fprint(w, content)
}

func playVideo(w http.ResponseWriter, r *http.Request) {
	video := GetUrlArg(r, "video")
	content := `<<body style="border:none;padding:0;margin:0;">
<div id="toolbox">
    <div>
        操作指南：<br/>
        f 切换全屏 <br/>
        方向键↑↓ 调整音量大小
        方向键←→ 改变视频进度
        as 调整视频播放速度
    </div>
    <!-- <div style="display: inline-block;width: 200px;">
        <label for="volume">Volume</label>
        <input type="range" id="volume" name="volume"
               min="0" max="11">
    </div> -->
</div>
<div style="top: calc(50% - 50px);
    left: calc(50% - 100px);
    position: absolute;
    z-index: 1000;
    background: #808080b3;
    width: 200px;
    height: 100px;
    border-radius: 8px;
    text-align: center;display: none;
    line-height: 100px;" id="tip">ibas</div>
<video controls="" autoplay="" name="media" style="width: 100%;" id="video">
    <source src="` + video + `">
</video>
<script>
    let video = null;
    let toolbox = null;
    let fs = 'fullscreen';
    let tip = null;
    const showTip = (label,value) => {
        tip.style.display = 'block';
        tip.innerText = label + value;
        setTimeout(() => {
            tip.style.display = 'none';
        },1000);
    };
    let tools = {
        fullscreen: {
            ele: false,
            key_event: {
                f() {
                    let fsV = video.getAttribute(fs);
                    if (fsV == 'true') {
                        toolbox.style.display = 'block';
                        document.exitFullscreen();
                        video.style.height = '';
                        video.setAttribute(fs,'false');
                    } else {
                        toolbox.style.display = 'none';
                        document.body.webkitRequestFullScreen();
                        video.style.height = '100%';
                        video.setAttribute(fs,'true');
                    }
                }
            }
        },
        volume: {
            ele: true, label: '音量',
            eleAttr: {
                id: 'volume', value: 100,
                min: 0,max: 100,
            },
            events: {
                onchange() {
                    video.volume = +this.value / 100;
                }
            },
            key_event: {
                up() {
                    let inp = document.getElementById('volume');
                    inp.value = +inp.value + 1;
                    inp.onchange();
                },
                down() {
                    let inp = document.getElementById('volume');
                    inp.value = +inp.value - 1;
                    inp.onchange();
                }
            }
        },
        time: {
            ele: true,
            label: '播放进度',
            eleAttr: {
                id: 'time', value: 0,
                min: 0,max() {
                    return new Promise(s => {
                        let id = setInterval(() => {
                            if (video.duration) {
                                s(video.duration);
                                clearInterval(id);
                            }
                        },1000);
                    });
                },
            },
            events: {
                onchange() {
                    video.currentTime = +this.value;
                }
            },
            key_event: {
                right() {
                    let inp = document.getElementById('time');
                    inp.value = +inp.value + video.duration / 100;
                    inp.onchange();
                },
                left() {
                    let inp = document.getElementById('time');
                    inp.value = +inp.value - video.duration / 100;
                    inp.onchange();
                }
            }
        },
        speed: {
            ele: true,
            label: '播放速度',
            eleAttr: {
                id: 'speed', value: 1,
                min: 0.5,max: 10, step: 0.2
            },
            events: {
                onchange() {
                    video.playbackRate = +this.value;
                }
            },
            key_event: {
                s() {
                    let inp = document.getElementById('speed');
                    let v = + inp.value + 0.2;
                    inp.value = v;
                    inp.onchange();
                },
                a() {
                    let inp = document.getElementById('speed');
                    inp.value = + inp.value - 0.2;
                    inp.onchange();
                }
            }
        }
    };
    const registerKeyEvent = function() {
        let rke = (new class{
            constructor() {
                this.events = {};
            }
            registerEvent(key,fn) {
                this.events[key] = fn;
            }
            // key = a~z 0~9 up\down\right\left
            emitKeyEvent(key) {
                if (key in this.events) {
                    this.events[key]();
                }
            }
        });
        document.onkeydown =  function keyDown(event){  // 方向键控制元素移动函数
            var event = event || window.event;  // 标准化事件对象
            switch(event.keyCode){  // 获取当前按下键盘键的编码
                case 37 :  // 按下左箭头键，向左移动5个像素
                    rke.emitKeyEvent("left");
                    break;
                case 39 :  // 按下右箭头键，向右移动5个像素
                    rke.emitKeyEvent("right");
                    break;
                case 38 :  // 按下上箭头键，向上移动5个像素
                    rke.emitKeyEvent("up");
                    break;
                case 40 :  // 按下下箭头键，向下移动5个像素
                    rke.emitKeyEvent("down");
                    break;
                default:
                    rke.emitKeyEvent(String.fromCharCode(event.keyCode).toLocaleLowerCase());
            }
            return false
        };
        return rke;
    };
    window.onload = function () {
        video = document.getElementById('video');
        toolbox = document.getElementById('toolbox');
        video.setAttribute(fs,'false');
        tip = document.getElementById('tip');
        let rke = registerKeyEvent();
        for (let i in tools) {
            if (tools[i].ele) {
                // 创建 dom
                // <div style="display: inline-block;width: 200px;">
                // <label for="volume">Volume</label>
                //         <input type="range" id="volume" name="volume"
                //     min="0" max="11">
                //         </div>
                let div = document.createElement('div');
                div.style.display = 'inline-block';
                div.style.width = '300px';
                let label = document.createElement('label');
                let input = document.createElement('input');
                let value = document.createElement('value');
                input.type = 'range';
                for (let j in tools[i].eleAttr) {
                    if (typeof tools[i].eleAttr[j] === 'function') {
                        tools[i].eleAttr[j]()
                        .then(v => {
                            input[j] = v;
                        });
                    } else {
                        input[j] = tools[i].eleAttr[j];
                    }
                }
                input.setAttribute('label',tools[i].label);
                value.innerText = tools[i].eleAttr.value;
                label.innerText = tools[i].label;
                if ('events' in tools[i]) {
                    for (let j in tools[i].events) {
                        if (j === 'onchange') {
                            input.onchange = function () {
                                tools[i].events[j].bind(this)();
                                value.innerText = this.value;
                                if (video.getAttribute(fs) === 'true') {
                                    showTip(this.getAttribute('label'),this.value);
                                }
                            }
                        } else {
                            input[j] = tools[i].events[j].bind(input);
                        }
                    }
                }
                div.append(label);div.append(input);div.append(value);toolbox.append(div);
            }
            if ('key_event' in tools[i]) {
                for (let j in tools[i].key_event) {
                    rke.registerEvent(j,tools[i].key_event[j]);
                }
            }
        }
    }
</script>
</body>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	fmt.Fprint(w, content)
}
