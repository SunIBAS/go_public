// const host = `http://${location.host}/api`;
const host = `http://localhost:8081/api`;

class PostParams {
    static tableObj = (itype,date,from,to,note) => {
        return {
            itype,from,to,date,note: note || '',num: ''
        }
    };
    constructor(method,content,ContentToString,parseMixin) {
        this.method = method || "";
        this.content = content || "";
        this.ContentToString = ContentToString || (_ => _);
        this.parseMixin = parseMixin || (_ => _);
    }
    // 理解为无需处理返回值的请求
    Post() {
        let that = this;
        return fetch(host, {
            method: 'POST', // *GET, POST, PUT, DELETE, etc.
            headers: {
                'Content-Type': 'application/json'
                // 'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: JSON.stringify({
                content: this.ContentToString(this.content),
                method: this.method
            }) // body data type must match "Content-Type" header
        }).then(_ => _.text()).then(_ => {
            return that.parseMixin(_);
        }).then(JSON.parse).catch(e => {
            alert(`请求[${that.method}]发生错误，${e.message}`);
            throw e;
        });
    }
    // 理解为需要将返回值做一定处理的请求
    PostParse() {
        return this.Post()
            .then(_ => {
                return JSON.parse(_.Content)
            })
    }

    DownloadFile(filename) {
        return fetch(host, {
            method: 'POST', // *GET, POST, PUT, DELETE, etc.
            headers: {
                'Content-Type': 'application/json'
                // 'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: JSON.stringify({
                content: this.ContentToString(this.content),
                method: this.method
            }) // body data type must match "Content-Type" header
        }).then(_ => _.blob()).then(blob => {
            let url = window.URL.createObjectURL(blob);
            let a = document.createElement('a');
            a.href = url;
            a.download = filename;
            a.click();
            window.URL.revokeObjectURL(url);
        })
    }
}

const apis = {
    // 获取目录节点的语言包
    // langObj = {id:'',en:'',zh:'',other:'',hsk:'',org:''}
    selectDI(date,from,itype) {
        return new PostParams("selectDI", JSON.stringify(PostParams.tableObj(itype,date,from))).Post()
    },
    select({itype,date,from,to,num,note,whos}) {
        let whereC = [];
        if (date) whereC.push(`date='${date}'`);
        if (from) whereC.push(`city_from='${from}'`);
        if (to) whereC.push(`city_to='${to}'`);
        if (num) whereC.push(`date='${num}'`);
        if (note) whereC.push(`date='${note}'`);
        if (whos) whereC.push(`whos='${whos}'`);
        if (whereC.length) {
            return new PostParams("select", JSON.stringify({
                table: itype,
                where: whereC.join(' and ')
            })).Post()
        } else {
            throw new Error("至少一个查询条件");
        }
    },
    insert({itype,date,from,to,num,note,whos}) {
        // param.From,param.To,param.Date,param.Num,param.Itype,param.Note
        return new PostParams("insert",JSON.stringify({itype,date,from,to,note,num,whos})).Post();
    },
    getNames() {
        return new PostParams("getNames").PostParse()
            .then(_ => {
                let names = [];
                _.forEach(n =>
                    names.push({value:n.name,label:n.name}));
                return names;
            });
    }
};

const orgUrl = [
    {
        "link": "http://www.nhc.gov.cn/cms-search/xxgk/searchList.htm?type=search",
        "title": "国家卫生健康委员会"
    },
    {
        "link": "http://wjw.beijing.gov.cn/xwzx_20031/wnxw/",
        "title": "北京卫生健康委员会"
    },
    {
        "link": "http://wsjk.gansu.gov.cn/channel/10910/index.html",
        "title": "甘肃卫生健康委员会"
    },
    {
        "link": "http://wsjkw.gd.gov.cn/zwyw_yqxx/",
        "title": "广东卫生健康委员会"
    },
    {
        "link": "http://wsjkw.henan.gov.cn/ztzl/xxgzbdfyyqfk/yqtb/",
        "title": "河南卫生健康委员会"
    },
    {
        "link": "http://wsjk.ln.gov.cn/wst_wsjskx/",
        "title": "辽宁卫生健康委员会"
    },
    {
        "link": "http://wsjkw.nx.gov.cn/yqfkdt/yqsd1.htm",
        "title": "宁夏回族自治区卫生健康委员会"
    },
    {
        "link": "http://wsjkw.shandong.gov.cn/ztzl/rdzt/qlzhxxgzbdfyyqfkgz/tzgg/",
        "title": "山东卫生健康委员会"
    },
    {
        "link": "https://wsjkw.sh.gov.cn/xwfb/index.html",
        "title": "上海卫生健康委员会（要细看）"
    },
    {
        "link": "https://wsjkw.zj.gov.cn/col/col1202194/index.html",
        "title": "浙江卫生健康委员会"
    },
    {
        "link": "http://wjw.shanxi.gov.cn/wjywl02/index.hrh",
        "title": "山西卫生健康委员会"
    },
    {
        "link": "http://ynswsjkw.yn.gov.cn/wjwWebsite/web/col?id=UU157976428326282067&cn=xxgzbd&pcn=ztlm&pid=UU145102906505319731",
        "title": "云南卫生健康委员会"
    },
    {
        "link": "http://wsjkw.gxzf.gov.cn/ztbd_49627/sszt/xxgzbdfyyqfk/yqtb/",
        "title": "广西壮族自治区卫生健康委员会"
    },
    {
        "link": "http://sxwjw.shaanxi.gov.cn/sy/wjyw/",
        "title": "陕西卫生健康委员会"
    },
    {
        "link": "http://wsjkw.sc.gov.cn/scwsjkw/gzbd01/ztwzlmgl.shtml",
        "title": "四川卫生健康委员会"
    },
    {
        "link": "http://wsjkw.hlj.gov.cn/pages/5df84bfaf6e9fa23e8848a48",
        "title": "黑龙江卫生健康委员会"
    },
    {
        "link": "http://wsjk.tj.gov.cn/ZTZL1/ZTZL750/YQFKZL9424/YQTB7440/",
        "title": "天津卫生健康委员会"
    },
    {
        "link": "http://wjw.fujian.gov.cn/xxgk/gzdt/wsjsyw/",
        "title": "福建省卫生健康委员会"
    },
    {
        "link": "http://www.hebwsjs.gov.cn/sjdt/index.jhtml",
        "title": "河北卫生健康委员会"
    },
    {
        "link": "http://hc.jiangxi.gov.cn/col/col38018/index.html",
        "title": "江西卫生健康委员会"
    },
    {
        "link": "http://wjw.jiangsu.gov.cn/col/col7290/index.html",
        "title": "江苏卫生健康委员会"
    },
    {
        "link": "http://wsjkw.cq.gov.cn/ztzl_242/qlzhxxgzbdfyyqfkgz/yqtb/",
        "title": "重庆卫生健康委员会"
    },
    {
        "link": "http://wjw.nmg.gov.cn/xwzx/xwfb/index.shtml",
        "title": "内蒙古卫生健康委员会"
    },
    {
        "link": "http://wsjkw.jl.gov.cn/xwzx/xwfb/",
        "title": "吉林卫生健康委员会"
    },
    {
        "link": "http://wjw.guizhou.gov.cn/xwzx_500663/zwyw/",
        "title": "贵州卫生健康委员会"
    },
    {
        "link": "http://wjw.ah.gov.cn/xwzx/gzdt/index.html",
        "title": "安徽卫生健康委员会"
    },
    {
        "link": "http://wjw.hubei.gov.cn/bmdt/ztzl/fkxxgzbdgrfyyq/xxfb/",
        "title": "湖北卫生健康委员会"
    },
    {
        "link": "http://wjw.hunan.gov.cn/wjw/xxgk/gzdt/zyxw_1/index.html",
        "title": "湖南卫生健康委员会"
    },
    {
        "link": "http://www.hainan.gov.cn/hainan/yqfkzzzzxtb/list_jjdyyqfkzzz.shtml",
        "title": "海南卫生健康委员会"
    },
    {
        "link": "http://wjw.xinjiang.gov.cn/hfpc/fkxxfyfkxx/fkxxfy_list.shtml",
        "title": "新疆维吾尔自治区卫生健康委员会"
    }
].map(_ => {
    return {
        ..._,time: 0
    }
});