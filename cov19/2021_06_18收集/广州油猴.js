/// 油猴
// ==UserScript==
// @name         20210618广东省wjw
// @namespace    http://tampermonkey.net/
// @version      0.1
// @description  try to take over the world!
// @author       You
// @match        http://*/*
// @icon         data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==
// @require      https://cdn.staticfile.org/jquery/3.5.0/jquery.min.js
// @grant        none
// ==/UserScript==

(function() {
    'use strict';
    function copy(text) {
        var input = document.createElement('input');
        input.setAttribute('readonly', 'readonly'); // 防止手机上弹出软键盘
        input.setAttribute('value', text);
        document.body.appendChild(input);
        // input.setSelectionRange(0, 9999);
        input.select();
        var res = document.execCommand('copy');
        document.body.removeChild(input);
        return res;
    }
    window.onload = function() {
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
        jQuery('.zoom_box').find('p,div').each((ind,p) => {
            // if (p.children.length === 0) return;
            if (p.innerText.indexOf('新增境外输入无症状感染者1、') !== -1) return;
            if (p.innerText.indexOf('新增境外输入无症状感染者1-') !== -1) return;
            if (p.innerText.indexOf('新增境外输入确诊病例1-') !== -1) return;
            if (p.innerText.indexOf('新增境外输入确诊病例1、') !== -1) return;
            if (p.innerText.indexOf('新增境外输入确诊病例2、') !== -1) return;
            if (p.innerText.indexOf('新增境外输入无症状感染者2、') !== -1) return;
            if (p.innerText.indexOf('新增境外输入无症状感染者3、') !== -1) return;
            if (p.innerText.indexOf('新增境外输入无症状感染者6、') !== -1) return;
            if (p.innerText.indexOf('新增境外输入无症状感染者4、') !== -1) return;
            let text = p.innerText.trim();
            if (text.startsWith('新增1例确诊病例情况') || text.startsWith('1例无症状感染者转确诊病例情况') || text.startsWith('确诊病例') || text.startsWith('新增确诊病例') || text.startsWith('新增1例境内确诊病例情况')) {
                境内确诊.push({
                    病例情况: '境内确诊',
                    年龄: parseInt(text.match(/[0-9]+岁/)[0]),
                    性别: text.indexOf('女') !== -1 ? '女' : '男',
                    是否无症状: '否',
                    是否疫苗: '-',
                    如何感染: text.indexOf('密切接触者') !== -1 ? '密切接触者' : "-",
                    哪个区: text.match(/广州市[\u4e00-\u9fa5]+区/)[0].substring(3)
                });
            } else if (text.startsWith('无症状感染者') || text.startsWith('新增无症状感染者')) {
                境内无症状感染者.push({
                    病例情况: '境内无症状感染者',
                    年龄: parseInt(text.match(/[0-9]+岁/)[0]),
                    性别: text.indexOf('女') !== -1 ? '女' : '男',
                    是否无症状: '是',
                    是否疫苗: '-',
                    如何感染: '-',
                    哪个区: 获取境外时间(text)
                });
            } else if (text.startsWith('新增1例境外输入无症状感染者') || text.startsWith('新增1例境外输入关联无症状感染者情况') || text.startsWith('境外输入无症状感染者') || text.startsWith('新增境外输入无症状感染者')) {
                境外输入无症状.push({
                    病例情况: '境外输入无症状',
                    年龄: parseInt(text.match(/[0-9]+岁/)[0]),
                    性别: text.indexOf('女') !== -1 ? '女' : '男',
                    是否无症状: '是',
                    是否疫苗: '-',
                    如何感染: '-',
                    哪个区: 获取境外时间(text)
                });
            } else if (text.startsWith('境外输入确诊病例') || text.startsWith('新增1例境外输入确诊病例情况') || text.startsWith('新增境外输入确诊病例') ) {
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
        });
        function dear每日一行数据() {
            let text = jQuery('.zoom_box')[0].innerText;
            let title = document.querySelector('body > div > div.main > div.cont_r.nr_cont_r > div.detailed_box > h4').innerText;
            let getNum = (reg) => {
                let ret = text.match(reg);
                if (ret) {
                    return ret[0].match(/[0-9]+/)[0];
                } else {
                    return '';
                }
            };
            let 日期 = title.match(/[0-9]+/g).join('.');
            let 每日新增境内 = getNum(/新增境内确诊病例[0-9]+例/);
            let 每日新增境外 = getNum(/新增境外输入确诊病例[0-9]+例/);
            let 累计确诊境内 = getNum(/境内确诊病例（含境外输入关联病例）[0-9]+例/);
            let 累计确诊境外 = getNum(/累计报告境外输入确诊病例[0-9]+例/);
            let 治疗人数 = getNum(/尚在院治疗[0-9]+例/);
            let 新增无症状 = getNum(/新增境外输入无症状感染者[0-9]+例/);
            let 出院 = getNum(/新增出院病例[0-9]+例/);
            let 累计出院 = getNum(/累计出院[0-9]+例/);
            console.log([日期,每日新增境内,每日新增境外,累计确诊境内,累计确诊境外,'',治疗人数,新增无症状,'',出院,累计出院].join('\t'));
            jQuery('body').append(jQuery(`<textarea style="
    position: absolute;
    top: 0;
    width: 100%;
    height: 40px;
">${[日期,每日新增境内,每日新增境外,累计确诊境内,累计确诊境外,'',治疗人数,新增无症状,'',出院,累计出院].join('\t')}</textarea>`));
        }
        function toText() {
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
            alert(`条目:${table.length}`);
            jQuery('body').append(jQuery(`<textarea style="
    position: absolute;
    top: 40px;
    width: 100%;
    height: 40px;
">${table.join('\r\n')}</textarea>`));
            //copy(table.join('\r\n'));
        }
        dear每日一行数据();
        toText();
    }
    // Your code here...
})();