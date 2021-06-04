该工具的用途是批量翻译 docx 和 xlsx 文件

总共分为四步完成该功能：1.格式化、2.翻译、3.更新、4.写出

第一步，格式化文件

第二步，翻译，右键翻译为中文，然后点击滚动，知道弹窗表示完成后，点击下载，然后点击下一页知道全部翻译完成

第三步，将第二步生成的翻译好的 json 文件写入数据库中

第四步，写出文件（翻译好的文件）

```text
1.格式化并生成文档模板
        将 inDir 下的所有文档格式化后写到 outDir 中
        将不递归遍历 inDir 目录
translateOffice -format inDir outDir

2.翻译，将开启一个服务，用于在浏览器上翻译文本
translateOffice -translate 8080

3.写回第二步生成的 json 文件
translateOffice -insert json

4.写出翻译好的文件
        需要提供第一步的 outDir 作为 base64Dir
        这里的 outDir 是翻译好的文档的输出文件夹
translateOffice -ok base64Dir outDir
```

第一步写出的 base64 文件的格式为

```xlsx + fileid + base64字符串``` 或
```docx + fileid + base64字符串```

fileid 长度为 36