// 为每日输入创建的 sql
let excel = `省份\t阿联酋\t土耳其\t加拿大\t刚果金\t几内亚\t埃塞俄比亚\t伊朗\t尼泊尔\t菲律宾
上海\t1\t1\t1\t1\t1\t2\t\t\t
广东\t\t\t\t\t\t\t1\t1\t
四川\t\t\t\t\t\t\t\t\t1
陕西\t\t\t\t\t\t\t\t\t1`.split('\n').map(_ => _.split('\t').map(_ => _.trim()));

let table = "daily_inj";
let date = "20200930";

for (let i = 1;i < excel.length;i++) {
    for (let j = 1;j < excel[0].length;j++) {
        if (excel[i][j]) {
            console.log(`insert into ${table}(city_from,city_to,date,num,note) values('${excel[0][j]}','${excel[i][0]}','${date}','${excel[i][j]}','excel');`);
        }
    }
}