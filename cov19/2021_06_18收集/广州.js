var lis = jQuery('.cont_list li');
var hrefs = [];
var fetchData = [];
lis.each((ind,li) => {
    if (ind) {
        hrefs.push({
            title: li.getElementsByTagName('a')[0].innerText,
            href: li.getElementsByTagName('a')[0].href
        });
    }
});

function dear每页数据(jqDom,titleText) {
    let 预估值 = 0;
    let 实际值 = 0;
    let 日期 = titleText.match(/[0-9]+/g).join('.');
    if (日期.length < 8) {
        日期 ='2020.' + 日期;
    }
    function dear每日一行数据() {
        let text = jqDom[0].innerText;
        let getNum = (reg) => {
            let ret = text.match(reg);
            if (ret) {
                return ret[0].match(/[0-9]+/)[0];
            } else {
                return '';
            }
        };

        let 每日新增境内 = getNum(/新增境内确诊病例[0-9]+例/) || 0;
        let 每日新增境外 = getNum(/新增境外输入确诊病例[0-9]+例/) || 0;
        let 累计确诊境内 = getNum(/境内确诊病例（含境外输入关联病例）[0-9]+例/);
        let 累计确诊境外 = getNum(/累计报告境外输入确诊病例[0-9]+例/);
        let 治疗人数 = getNum(/尚在院治疗[0-9]+例/);
        let 新增无症状 = getNum(/新增境外输入无症状感染者[0-9]+例/) || 0;
        let 出院 = getNum(/新增出院病例[0-9]+例/);
        let 累计出院 = getNum(/累计出院[0-9]+例/);
        console.log([日期,每日新增境内,每日新增境外,累计确诊境内,累计确诊境外,'',治疗人数,新增无症状,'',出院,累计出院].join('\t'));
//         jQuery('body').append(jQuery(`<textarea style="
//     position: absolute;
//     top: 0;
//     width: 100%;
//     height: 40px;
// ">${[日期,每日新增境内,每日新增境外,累计确诊境内,累计确诊境外,'',治疗人数,新增无症状,'',出院,累计出院].join('\t')}</textarea>`));
        预估值 = (+每日新增境内) + (+每日新增境外) + (+新增无症状);
        return [日期,每日新增境内,每日新增境外,累计确诊境内,累计确诊境外,'',治疗人数,新增无症状,'',出院,累计出院].join('\t');
    }
    function toText() {
        let 境内确诊 = [
            // {
            //     病例情况: '境内确诊',
            //     年龄: '45',
            //     性别: '',
            // }
        ];
        let 境外输入无症状 = [];
        let 境外输入确诊病例 = [];
        let 境内无症状感染者 = [];
        function 获取境外时间(text) {
            let find = text.indexOf('从');
            if (find === -1) {
                return '';
            } else {
                while (find && text[find] !== '。') {
                    find--;
                }
                return text.substring(find + 1);
            }
        }
        let ps = jqDom.find('p,div');
        let dears = {
            境内确诊(text) {
                境内确诊.push({
                    病例情况: '境内确诊',
                    年龄: parseInt(text.match(/[0-9]+岁/)[0]),
                    性别: text.indexOf('女') !== -1 ? '女' : '男',
                    是否无症状: '否',
                    是否疫苗: '-',
                    如何感染: text.indexOf('密切接触者') !== -1 ? '密切接触者' : "-",
                    哪个区: text.match(/广州市[\u4e00-\u9fa5]+区/)[0].substring(3)
                });
            },
            境内无症状感染者(text) {
                境内无症状感染者.push({
                    病例情况: '境内无症状感染者',
                    年龄: parseInt(text.match(/[0-9]+岁/)[0]),
                    性别: text.indexOf('女') !== -1 ? '女' : '男',
                    是否无症状: '是',
                    是否疫苗: '-',
                    如何感染: '-',
                    哪个区: 获取境外时间(text)
                });
            },
            境外输入无症状(text) {
                境外输入无症状.push({
                    病例情况: '境外输入无症状',
                    年龄: parseInt(text.match(/[0-9]+岁/)[0]),
                    性别: text.indexOf('女') !== -1 ? '女' : '男',
                    是否无症状: '是',
                    是否疫苗: '-',
                    如何感染: '-',
                    哪个区: 获取境外时间(text)
                });
            },
            境外输入确诊病例(text) {
                境外输入确诊病例.push({
                    病例情况: '境外输入确诊病例',
                    年龄: parseInt(text.match(/[0-9]+岁/)[0]),
                    性别: text.indexOf('女') !== -1 ? '女' : '男',
                    是否无症状: '否',
                    是否疫苗: '-',
                    如何感染: '-',
                    哪个区: 获取境外时间(text)
                });
            }
        };
        ps.each((ind,p) => {
            // if (p.children.length === 0) return;
            if (/新增境外输入无症状感染者[0-9]+[、-]{1}/.test(p.innerText)) return;
            if (/新增境外输入确诊病例[0-9]+[、-]{1}/.test(p.innerText)) return;
            if (/新增境外输入无症状感染者[0-9]+[、-]{1}/.test(p.innerText)) return;
            // if (p.innerText.indexOf('新增境外输入无症状感染者1') !== -1) return;
            // if (p.innerText.indexOf('新增境外输入确诊病例1-') !== -1) return;
            // if (p.innerText.indexOf('新增境外输入确诊病例1、') !== -1) return;
            // if (p.innerText.indexOf('新增境外输入确诊病例2、') !== -1) return;
            // if (p.innerText.indexOf('新增境外输入无症状感染者2、') !== -1) return;
            // if (p.innerText.indexOf('新增境外输入无症状感染者3、') !== -1) return;
            // if (p.innerText.indexOf('新增境外输入无症状感染者6、') !== -1) return;
            // if (p.innerText.indexOf('新增境外输入无症状感染者4、') !== -1) return;
            // if (p.innerText.indexOf('新增境外输入无症状感染者5、') !== -1) return;
            let text = p.innerText.trim();
            if (text.startsWith('新增1例确诊病例情况') || text.startsWith('1例无症状感染者转确诊病例情况') || text.startsWith('确诊病例') || text.startsWith('新增确诊病例') || text.startsWith('新增1例境内确诊病例情况')) {
                if (text.indexOf('岁') !== -1) {
                    dears.境内确诊(text);
                } else {
                    dears.境内确诊(ps[ind + 1].innerText.trim());
                }
            } else if (text.startsWith('无症状感染者') || text.startsWith('新增无症状感染者')) {
                if (text.indexOf('岁') !== -1) {
                    dears.境内无症状感染者(text);
                } else {
                    dears.境内无症状感染者(ps[ind + 1].innerText.trim());
                }
            } else if (text.startsWith('新增1例境外输入无症状感染者') || text.startsWith('新增1例境外输入关联无症状感染者情况') || text.startsWith('境外输入无症状感染者') || text.startsWith('新增境外输入无症状感染者')) {
                if (text.indexOf('岁') !== -1) {
                    dears.境外输入无症状(text);
                } else {
                    dears.境外输入无症状(ps[ind + 1].innerText.trim());
                }
            } else if (text.startsWith('境外输入确诊病例') || text.startsWith('新增1例境外输入确诊病例情况') || text.startsWith('新增境外输入确诊病例') ) {
                if (text.indexOf('岁') !== -1) {
                    dears.境外输入确诊病例(text);
                } else {
                    dears.境外输入确诊病例(ps[ind + 1].innerText.trim());
                }
            }
        });
        let table = [];
        function _toText(d) {
            let dd = [d.病例情况,d.年龄,d.性别,d.是否无症状,d.是否疫苗,d.如何感染,d.哪个区];
            // console.log();
            table.push(dd.join('\t'));
        };
        境内确诊.forEach(_toText);
        境内无症状感染者.forEach(_toText);
        境外输入确诊病例.forEach(_toText);
        境外输入无症状.forEach(_toText);
        console.log(table.length);
        console.log(table.join('\r\n'));
        // alert(`条目:${table.length}`);
//         jQuery('body').append(jQuery(`<textarea style="
//     position: absolute;
//     top: 40px;
//     width: 100%;
//     height: 40px;
// ">${table.join('\r\n')}</textarea>`));
        实际值 = table.length;
        return table.join('\r\n');
        //copy(table.join('\r\n'));
    }
    return [
        dear每日一行数据(),
        日期.startsWith('2020') ? "" : toText(),
        预估值,实际值
    ];
}
var funDownload = function (content, filename) {
    // 创建隐藏的可下载链接
    var eleLink = document.createElement('a');
    eleLink.download = filename;
    eleLink.style.display = 'none';
    // 字符内容转变成blob地址
    var blob = new Blob([content]);
    eleLink.href = URL.createObjectURL(blob);
    // 触发点击
    document.body.appendChild(eleLink);
    eleLink.click();
    // 然后移除
    document.body.removeChild(eleLink);
};
var allContent = "";

function fetchAll() {
    let d = hrefs.shift();
    fetchData.unshift({
        ...d,
        text: '',
    });
    if (!/[0-9]+月[0-9]+日/.test(fetchData[0].title)) {
        if (hrefs.length) {
            fetchAll();
        } else {
            funDownload(allContent,"a.txt");
            console.log('请求结束');
        }
        return ;
    }
    console.log(`开始处理${fetchData[0].title}`);
    fetch(d.href,{method: 'get'}).then(_ => _.text()).then(data =>{
        let jqDom = jQuery(data).find('.zoom_box');
        let ret = dear每页数据(jqDom,jQuery(data).find('.detailed_box h4')[0].innerText);
        allContent += `=====================================
${ret[2] === ret[3] ? '正常' : "**************不正常*************"}
${ret[0]}
${ret[1]}
`;
        if (hrefs.length) {
            fetchAll();
        } else {
            funDownload(allContent,"a.txt");
            console.log('请求结束');
        }
    });
}
fetchAll();

