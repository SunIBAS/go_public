const datas = require('./out');

// let o = { Id: '13',
//     Content:
//         '{"status":"success","data":{"buy_price":"46971.93","sell_price":"46851.1"},"mobile":null}',
//     Time: '2021-02-23 22:43:55' };
let option = {
    tooltip: {
        trigger: 'axis'
    },
    legend: {
        data: ['买价','卖价']
    },
    toolbox: {
        feature: {
            saveAsImage: {}
        }
    },
    xAxis: {
        type: 'category',
        data: []
    },
    yAxis: {
        type: 'value'
    },
    series: [{
        name: '买价',
        data: [],
        type: 'line',
        smooth: true
    },{
        name: '卖价',
        data: [],
        type: 'line',
        smooth: true
    }]
};
let count = 0;
let max = -1000000;
let min = 1000000;
datas.reverse().forEach(o => {
    if (count > 1e20) {
        return;
    }
    try {
        o.Content = JSON.parse(o.Content);
        option.xAxis.data.push(o.Time);
        option.series[0].data.push(o.Content.data.buy_price);
        option.series[1].data.push(o.Content.data.sell_price);
        o.Content.data.buy_price = +o.Content.data.buy_price;
        o.Content.data.sell_price = +o.Content.data.sell_price;
        max = max > o.Content.data.buy_price ? max : o.Content.data.buy_price;
        max = max > o.Content.data.sell_price ? max : o.Content.data.sell_price;
        min = min < o.Content.data.buy_price ? min : o.Content.data.buy_price;
        min = min < o.Content.data.sell_price ? min : o.Content.data.sell_price;
        count++;
    } catch (e) {
    }
});
option.xAxis.data.reverse();
option.series[0].data.reverse();
option.series[1].data.reverse();
console.log(`option = {
    tooltip: {
        trigger: 'axis'
    },
    legend: {
        data: ['买价','卖价']
    },
    toolbox: {
        feature: {
            saveAsImage: {}
        }
    },
    dataZoom: [
        {
            type: 'inside',
            xAxisIndex: [0, 1],
            start: 10,
            end: 100
        },
        {
            show: true,
            xAxisIndex: [0, 1],
            type: 'slider',
            bottom: 10,
            start: 10,
            end: 100
        }
    ],
    xAxis: {
        type: 'category',
        data: ${JSON.stringify(option.xAxis.data)}
    },
    yAxis: {
        type: 'value',
        min: ${min},
        max: ${max},
    },
    series: [{
        name: '买价',
        data: ${JSON.stringify(option.series[0].data)},
        type: 'line',
        smooth: true
    },{
        name: '卖价',
        data: ${JSON.stringify(option.series[1].data)},
        type: 'line',
        smooth: true
    }]
}`);
