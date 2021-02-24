package main

import (
	uuid "github.com/satori/go.uuid"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"public.sunibas.cn/go_public/utils/MyRegex"
	"public.sunibas.cn/go_public/utils/ParseOffice"
	"strings"
)

func main() {

}

// 格式化
// docFilePath doc 文档的全路径
// tmpPath 	   临时文件路径
// htmlPath    用于翻译的html
func normalWord(docFilePath, tmpPath, htmlPath string) {
	nor := ParseOffice.RpWord{
		InFile:  docFilePath,
		OutFile: tmpPath,
		Merge:   true,
	}
	// 初始化
	nor.Init()
	// 要求合并 xml 节点（并提取所有的文件同时进行去重）
	nor.MergeT()
	// 对每个句子生成唯一标签
	kvs := buildKV(nor.Words)
	nor.Replace(kvs)
	// 写出模板文件
	nor.WriteOut()
	writeOutAsHTML(kvs, htmlPath)
}
func getUUID() string {
	var err error
	return uuid.Must(uuid.NewV4(), err).String()
}

// 生成 kv，主要是去重和给定标志
func buildKV(lines []string) map[string]string {
	kvs := map[string]string{}
	id := ""
	for _, rawline := range lines {
		line := strings.Replace(rawline, "'", "''", -1)
		if MyRegex.TestNoWord(line) {
			continue
		}
		if _, ok := kvs[line]; !ok {
			id = getUUID()
			kvs[line] = id
		}
	}
	return kvs
}

// 写出用于翻译的 html 文件
func writeOutAsHTML(kvs map[string]string, htmlPath string) {
	lines := []string{}
	for k, v := range kvs {
		lines = append(lines, `<div id='`+k+`'>`+v+`</div`)
	}
	html := strings.Replace(TranslateHTML, "$content$", strings.Join(lines, ""), 1)
	DirAndFile.WriteWithIOUtil(htmlPath, html)
}
