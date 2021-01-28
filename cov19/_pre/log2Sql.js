const fs = require('fs');

let text = fs.readFileSync('log.txt','utf-8').split('\n')
    .map(_ => _.substring(20).trim());

let sqls = [];

let dateReg = /[0-9]+月[0-9]+日/;
let int2str2 = n => (n.length < 2 ? '0' : '') + n;

function class_1(ind) {
    let table = text[ind];
    let m = text[ind + 1].match(dateReg);
    let md = "";
    if (m && m.length) {
        md = m[0].match(/[0-9]+/g);
        md[0] = int2str2(md[0]);
        md[1] = int2str2(md[1]);
        if (+md[0] > 2) {
            md = `2020${md[0]}${md[1]}`;
        } else {
            md = `2021${md[0]}${md[1]}`;
        }
    } else {
        console.log(text[ind]);
        md = text[ind + 1];
    }
    for (let i = ind + 2;i < text.length;i++) {
        if (text[i] === "--end--") {
            return i;
        } else {
            let c = text[i].split('\t').map(_ => _.trim()).join("','");
            sqls.push(`insert into ${table}(city_from,city_to,num,date,note) values('${c}','${md}','excel');`);
        }
    }
}

for (let i = 0;i < text.length;i++) {
    if (text[i] === "--start--") {
        i = class_1(i + 1);
    }
}

fs.writeFileSync('sqls.sql',sqls.join('\r\n'));