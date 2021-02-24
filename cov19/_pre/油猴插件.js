// ==UserScript==
// @name         广东卫生健康委员会
// @namespace    http://tampermonkey.net/
// @version      0.1
// @description  try to take over the world!
// @author       You
// @match        *.gov.cn/*
// @grant        none
// ==/UserScript==

(function() {
    'use strict';
    let colorsFn = {
        red: s => `<span style="color: #f4369f;background: #8BC34A;font-size: 25px;font-weight: bold;">${s}</span>`,
        blue: s => `<span style="color:blue;background: #8BC34A;font-size: 25px;font-weight: bold;">${s}</span>`,
        green: s => `<span style="color:green;background: #8bc34a63;font-size: 25px;font-weight: bold;">${s}</span>`,
        color1: s => `<span style="color: #80005f;background: #4ac36a9e;font-size: 25px;font-weight: bold;">${s}</span>`,
        color2: s => `<span style="color: #b33926;background: #b4c5bee0;font-size: 25px;font-weight: bold;">${s}</span>`,
        hide: s => `<span style="color: #96a766;background: #abe864;font-size: 25px;font-weight: bold;">${s}</span>`,
    };
    function replaceStr(str,reg,colorFn) {
        let m = str.match(reg);
        if (m && m.length) {
            return str.substring(0,m.index) + colorFn(m[0]) + replaceStr(str.substring(m.index + m[0].length),reg,colorFn);
        } else {
            return str;
        }
    }
    function getPs(selector,pName = 'p') {
        let p = document.querySelector(selector);
        if (p) {
            return p.getElementsByTagName(pName || 'p');
        } else {
            return [];
        }
    }
    function psWithRulues(ps,rules) {
        for (let i = 0;i < ps.length;i++) {
            let str = ps[i].innerText;
            rules.forEach(rr => {
                str = replaceStr(str,rr[0],rr[1]);
            });
            ps[i].innerHTML = str;
        }
    }
    let run = false;
    function load() {
        if (run) return;
        run = true;
        if (location.href.includes('http://wsjkw.gd.gov.cn/zwyw_yqxx/content/')) { // 广东卫生健康委员会
            setInterval(function() {
                let a = document.getElementById('webIco');
                if (a) a.remove();
            },2000);
            let ps = getPs(".content-content");//document.querySelector('.con_font').getElementsByTagName('p');
            let repRules = [
                [/全省新增境外输入确诊病例[0-9]+例/,colorsFn.red],
                [/新增[0-9]+例境外输入确诊病例/,colorsFn.red],
                [/新增[0-9]+例本土无症状感染者/,colorsFn.blue],
                [/新增境外输入无症状感染者[0-9]+例/,colorsFn.blue],
                [/新增出院[0-9]+例/,colorsFn.green],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://www.nhc.gov.cn/cms-search/xxgk/')) { // 全国
            let ps = getPs(".con_font");//document.querySelector('.con_font').getElementsByTagName('p');
            let repRules = [
                [/和新疆生产建设兵团报告新增确诊病例[0-9]+例/,colorsFn.red],
                [/其中境外输入病例[0-9]+例/,colorsFn.blue],
                [/[本土|本地]*病例[0-9]+例/,colorsFn.green],
                [/和新疆生产建设兵团报告新增无症状感染者[0-9]+例/,colorsFn.green],
                [/境外输入[0-9]+例/,colorsFn.green],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wjw.beijing.gov.cn/xwzx')) { // 北京
            let ps = getPs("#zoom > div");
            let repRules = [
                [/北京市新增[0-9]+例[本土|本地]*新冠肺炎确诊病例/,colorsFn.red],
                [/新增[0-9]+例[本土|本地]*确诊病例/,colorsFn.red],
                [/新增[0-9]+例境外输入*确诊病例/,colorsFn.red],
                [/[0-9]+例无症状感染者/,colorsFn.blue],
                [/新增[0-9]+例境外输入无症状感染者/,colorsFn.green],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wjw.shanxi.gov.cn/')) { // 山西
            let ps = getPs("#container > div.cont > div.page-right > div.zeTop > div.boxC > div.ze-art");
            let repRules = [
                [/山西省[本土|本地]*新增新冠肺炎确诊病例[0-9]+例/,colorsFn.red],
                [/山西无新增境外输入确诊病例/,colorsFn.blue],
                [/新增境外输入确诊病例[0-9]+例/,colorsFn.blue],
                [/山西省无新增无症状感染者/,colorsFn.green],
                [/新增无症状感染者[0-9]+例/,colorsFn.green],
                [/无新增/,colorsFn.green],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wsjkw.sc.gov.cn/')) { // 四川
            let ps = getPs("body > div.wy_dis_wiap > div.wy_dis_main > div.wy_contMain.fontSt");
            let repRules = [
                [/无新增/,colorsFn.red],
                [/新增新型冠状病毒肺炎确诊病例[0-9]+例/,colorsFn.green],
                [/新增无症状感染者[0-9]+例/,colorsFn.blue],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wsjk.gansu.gov.cn/')) { // 甘肃
            let ps = getPs("#contents");
            let repRules = [
                [/无新增/,colorsFn.red],
                [/新增新型冠状病毒肺炎确诊病例[0-9]+例/,colorsFn.green],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wsjkw.gxzf.gov.cn/')) { // 广西
            let ps = getPs("body > div.wrap.bg-white > div.more > div > div.article-con > div.view.TRS_UEDITOR.trs_paper_default.trs_web");
            if (!ps.length) {
                ps = getPs("body > div.wrap.bg-white > div.more > div > div.article-con > div.view.TRS_UEDITOR.trs_paper_default");
            }
            let repRules = [
                [/无新增/,colorsFn.red],
                [/新增[0-9]+例境外输入无症状感染者/,colorsFn.red],
                [/现有/,colorsFn.color2],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://ynswsjkw.yn.gov.cn/')) { // 云南
            let ps = getPs("#content");
            let repRules = [
                [/无新增/,colorsFn.red],
                [/输入无症状感染者[0-9]+例/,colorsFn.blue],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://sxwjw.shaanxi.gov.cn/')) { // 陕西
            let ps = getPs("body > div.w-content-bg > div > div > div.message-box > div.clearfix.news-detail > div > div.view.TRS_UEDITOR.trs_paper_default.trs_web");
            let repRules = [
                [/无新增/,colorsFn.red],
                [/新增报告[0-9]+例本土确诊病例/,colorsFn.blue],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('https://wsjkw.sh.gov.cn/')) { // 上海
            let ps = getPs("#ivs_content");
            let repRules = [
                [/上海新增[0-9]+例[本土|本地]*确诊病例/,colorsFn.red],
                [/无新增/,colorsFn.blue],
                [/病例[0-9]+/,colorsFn.green],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wsjkw.nx.gov.cn/')) { // 宁夏
            let ps = getPs("#vsb_content_4");
            let repRules = [
                [/全区无新增新冠肺炎确诊病例，无新增疑似病例，无无症状感染者报告/,colorsFn.red],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wsjkw.shandong.gov.cn/')) { // 山东
            let ps = getPs("body > div.liebiaonei > div > div.text > div.view.TRS_UEDITOR.trs_paper_default.trs_web","section");
            let repRules = [
                [/无新增/,colorsFn.red],
                [/境外输入无症状感染者[0-9]+例/,colorsFn.red],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wsjk.ln.gov.cn/')) { // 辽宁
            let ps = getPs("#ContTextSize > div > div");
            let repRules = [
                [/无新增新冠肺炎确诊病例/,colorsFn.red],
                [/无新增无症状感染者/,colorsFn.blue],
                [/新增本土确诊病例治愈出院[0-9]+例/,colorsFn.hide],
                [/新增[0-9]+例境外输入新冠肺炎确诊病例/,colorsFn.green],
                [/辽宁省新增境外输入新冠肺炎确诊病例[0-9]+例/,colorsFn.green],
                [/境外输入无症状感染者[0-9]+例/,colorsFn.green],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wsjkw.henan.gov.cn/2021')) { // 河南
            let ps = getPs("#artibody");
            let repRules = [
                [/我省无新增确诊病例、疑似病例、无症状感染者/,colorsFn.red],
                [/我省无新增确诊病例和疑似病例/,colorsFn.green],
                [/新增无症状感染者[0-9]+例/,colorsFn.blue],
                [/为境外输入）/,colorsFn.color1],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('https://wsjkw.zj.gov.cn/')) { // 浙江
            let ps = getPs("#zoom");
            let repRules = [
                [/新增确诊病例[0-9]+例/,colorsFn.red],
                [/新增无症状感染者[0-9]+例/,colorsFn.blue],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wsjkw.hlj.gov.cn/')) { // 黑龙江
            let ps = getPs("body > div.xxgk_bg > div.xxgkb_content > div.xxgk_bottom > div.gknb_bottom_right > div.gknb_content");
            let repRules = [
                [/黑龙江省新增新冠肺炎确诊病例[0-9]+例/,colorsFn.red],
                [/新增无症状感染者[0-9]+例/,colorsFn.blue],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wsjk.tj.gov.cn/')) { // 天津
            let ps = getPs("#zoom > div.view.TRS_UEDITOR.trs_paper_default.trs_web","section");
            if (!ps.length) {
                ps = getPs("#zoom > div.view.TRS_UEDITOR.trs_paper_default.trs_web","p");
            }
            let repRules = [
                [/天津市无新增[本土|本地]*新冠肺炎确诊病例/,colorsFn.red],
                [/新增[0-9]+例境外输入性新冠肺炎确诊病例/,colorsFn.blue],
                [/新增[0-9]+例无症状感染者/,colorsFn.blue],
                [/新增无症状感染者[0-9]+例/,colorsFn.red],
                [/无新增/,colorsFn.green],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://www.hebwsjs.gov.cn/')) { // 河北
            let ps = getPs("#zoom");
            let repRules = [
                [/河北省新增[0-9]+例[本土|本地]*新型冠状病毒肺炎确诊病例/,colorsFn.red],
                [/河北省新增[0-9]+例[本土|本地]*无症状感染者/,colorsFn.blue],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wsjkw.jl.gov.cn/')) { // 吉林
            let ps = getPs("#zoom > div > div");
            if (!ps.length) {
                ps = getPs("#zoom > div > div",'div');
            }
            if (!ps.length) {
                run = false;
                setTimeout(load,2000);
                return;
            }
            if (ps[0].parentElement.childNodes[0].textContent.trim()) {
                let p = document.createElement('p');
                p.innerHTML = ps[0].parentElement.childNodes[0].textContent;
                ps[0].parentElement.childNodes[0].remove();
                ps[0].parentElement.insertBefore(p,ps[0]);
                let pps = [p];
                for (let i = 0;i < ps.length;i++) {
                    pps.push(ps[i]);
                }
                ps = pps;
            }
            let repRules = [
                [/全省新增[本土|本地]*确诊病例[0-9]+例/,colorsFn.red],
                [/全省新发现确诊病例[0-9]+例/,colorsFn.red],
                [/有[0-9]+例转为确诊病例/,colorsFn.red],
                [/新增无症状感染者[0-9]+例/,colorsFn.blue],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wjw.jiangsu.gov.cn/')) { // 江苏
            let ps = getPs("#barrierfree_container > div.w1100.center > div.main-fl.bt-left");
            let repRules = [
                [/新增境外输入确诊病例[0-9]+例/,colorsFn.red],
                [/当日新增境外输入无症状感染者[0-9]+例/,colorsFn.red],
                [/当日新增无症状感染者[0-9]+例/,colorsFn.blue],
                [/无新增/,colorsFn.red],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wjw.nmg.gov.cn/')) { // 内蒙古
            let ps = getPs("#Zoom");
            let repRules = [
                [/无新增/,colorsFn.red],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wjw.hubei.gov.cn/')) { // 湖北
            let ps = getPs("#article-box > div");
            let repRules = [
                [/增新冠肺炎确诊病例[0-9]+例/,colorsFn.red],
                [/新增疑似病例[0-9]+例/,colorsFn.blue],
                [/新增无症状感染者[0-9]+例/,colorsFn.green],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wjw.guizhou.gov.cn/')) { // 贵州
            let ps = getPs("#Zoom");
            let repRules = [
                [/无新增/,colorsFn.red],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wjw.ah.gov.cn/')) { // 安徽
            let ps = getPs("body > div > div.container > div > div.ls-content > div > div.wzcon.j-fontContent.clearfix","p");
            let repRules = [
                [/无新增/,colorsFn.red],
                [/新增境外输入确诊病例[0-9]+例/,colorsFn.blue],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wsjkw.cq.gov.cn/')) { // 重庆
            let ps = getPs("body > div.main > div > div > div > div.zwxl-article > div.view.TRS_UEDITOR.trs_paper_default.trs_word");
            let repRules = [
                [/无新增新冠肺炎确诊病例报告/,colorsFn.red],
                [/无新增境外输入新冠肺炎确诊病例报告/,colorsFn.blue],
                [/报告新增无症状感染者[0-9]+例/,colorsFn.red],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wjw.hunan.gov.cn/')) { // 湖南
            let ps = getPs("#j-show-body",'section');
            let repRules = [
                [/湖南省报告新增新型冠状病毒肺炎确诊病例[0-9]例/,colorsFn.red],
                [/湖南省当日新增新型冠状病毒肺炎无症状感染者[0-9]+例/,colorsFn.blue],
                [/湖南省报告新增新型冠状病毒肺炎无症状感染者[0-9]+例/,colorsFn.blue],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://wjw.fujian.gov.cn/')) { // 福建
            let ps = getPs("#detailContent > div");
            let repRules = [
                [/新增境外输入无症状感染者[0-9]+例/,colorsFn.red],
                [/新增境外输入确诊病例/,colorsFn.blue],
                [/新增本地疑似病例/,colorsFn.green],
                [/新增[本地|本土]*无症状感染者/,colorsFn.color1],
            ];
            psWithRulues(ps,repRules);
        } else if (location.href.includes('http://hc.jiangxi.gov.cn/')) { // 江西
            let ps = getPs("#zoom");
            let repRules = [
                [/无新增/,colorsFn.red],
            ];
            psWithRulues(ps,repRules);
        }
    }
    window.onload = load;
    setTimeout(load,2000);
    // Your code here...
})();