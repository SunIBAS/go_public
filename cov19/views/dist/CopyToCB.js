function mmcopy(txt) {
    if (document.execCommand("copy")) {
        const input = document.createElement("input"); // 创建一个新input标签;
        input.setAttribute("readonly", "readonly"); // 设置input标签只读属性
        input.setAttribute("value", txt); // 设置input value值为需要复制的内容
        document.body.appendChild(input); // 添加input标签到页面
        input.select(); // 选中input内容
        input.setSelectionRange(0, 9999); // 设置选中input内容范围
        document.execCommand("copy"); // 复制
        document.body.removeChild(input);  // 删除新创建的input标签
    }
}