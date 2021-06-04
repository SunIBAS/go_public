package NormalFile

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"os"
	"path"
	"public.sunibas.cn/go_public/translateOffice/Server"
	"public.sunibas.cn/go_public/translateOffice/db"
	"public.sunibas.cn/go_public/utils/Console"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"public.sunibas.cn/go_public/utils/MyRegex"
	"public.sunibas.cn/go_public/utils/ParseOffice"
	"strconv"
	"strings"
)

var LangDatabase *db.LangDb
var DirDatabase *db.DirDb
var FileDatabase *db.FileDb
var Tag string
var tmpFilePath = "__tmp__file__.tmp"
var idLength = 36

func NormalFile(
	_lang_db_ *db.LangDb,
	_dir_db_ *db.DirDb,
	_file_db_ *db.FileDb,
	inDir,
	outDir,
	tag string) {
	LangDatabase = _lang_db_
	DirDatabase = _dir_db_
	FileDatabase = _file_db_
	//LangDatabase.Clear()
	Tag = tag
	total := 0
	// 设置保存 base64 文件的方法
	if _dir_db_ != nil {
		dearDir(inDir)
	} else {
		files := DirAndFile.GetSubDirOrFile(inDir)
		for _, f := range files {
			if f.File {
				if f.Name[0] == '~' {
					continue
				}
				if len(f.Name) > 5 && strings.ToLower(f.Name[len(f.Name)-4:]) == db.DOCX {
					fmt.Println("正在处理文件 -> " + f.Name)
					total += normalWord(path.Join(inDir, f.Name), path.Join(outDir, f.Name[0:len(f.Name)-5]+".base64"))
				} else if len(f.Name) > 5 && strings.ToLower(f.Name[len(f.Name)-4:]) == db.XLSX {
					fmt.Println("正在处理文件 -> " + f.Name)
					total += normalExcel(path.Join(inDir, f.Name), path.Join(outDir, f.Name[0:len(f.Name)-5]+".base64"))
				}
			}
		}
	}

	fmt.Println("总共插入 " + strconv.Itoa(total) + " 条记录")
}

func normalExcel(xlsxFilePath, base64FilePath string) int {
	count := 0
	fileId := GetUUID()
	re := ParseOffice.RpExcel{
		InFile:  xlsxFilePath,
		OutFile: tmpFilePath,
		Merge:   true,
	}
	// 初始化
	re.Init()
	// 要求合并 xml 节点（并提取所有的文件同时进行去重）
	re.MergeT()
	// 将所有文字存到数据库中并为文字生成标记符号
	wbKvs := buildKV(re.Words)
	count += saveKV(fileId, db.WorkBook, wbKvs)
	snKvs := buildKV(re.GetSheetName())
	count += saveKV(fileId, db.SheetName, snKvs)
	for k, v := range wbKvs {
		wbKvs[k] = "{{" + v + "}}"
	}
	re.Replace(wbKvs)
	re.ReplaceSheetNameToUUID(snKvs)
	re.WriteOut()
	DirAndFile.WriteWithIOUtil(base64FilePath, db.XLSX+fileId+DirAndFile.ReadFileAsBase64String(re.OutFile))
	os.Remove(re.OutFile)
	return count
}

// 格式化
// docFilePath 		doc 文档的全路径
// base64FilePath	模板文件将以 base64 保存，并且前 idLength 为文件的 id
func normalWord(docFilePath, base64FilePath string) int {
	nor := ParseOffice.RpWord{
		InFile:  docFilePath,
		OutFile: tmpFilePath,
		Merge:   true,
	}
	fileId := GetUUID()
	// 初始化
	nor.Init()
	// 要求合并 xml 节点（并提取所有的文件同时进行去重）
	nor.MergeT()
	// 对每个句子生成唯一标签
	kvs := buildKV(nor.Words)
	count := saveKV(fileId, db.Word, kvs)
	for k, v := range kvs {
		kvs[k] = "{{" + v + "}}"
	}
	nor.Replace(kvs)
	// 写出模板文件
	nor.WriteOut()
	//writeOutAsHTML(kvs, htmlPath)
	DirAndFile.WriteWithIOUtil(base64FilePath, db.DOCX+fileId+DirAndFile.ReadFileAsBase64String(nor.OutFile))
	// 删除临时文件
	os.Remove(nor.OutFile)
	return count
}

func GetUUID() string {
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
			id = GetUUID()
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
	html := strings.Replace(Server.TranslateHTML, "$content$", strings.Join(lines, ""), 1)
	DirAndFile.WriteWithIOUtil(htmlPath, html)
}

func saveKV(fileid string, type_ string, kv map[string]string) int {
	count := 0
	values := []string{}
	keys := []string{}
	for key, value := range kv {
		values = append(values, value)
		keys = append(keys, key)
		count++
		if len(values) >= 500 {
			LangDatabase.InsertValues(values, fileid, type_, keys)
			values = []string{}
			keys = []string{}
		}
	}
	if len(values) > 0 {
		LangDatabase.InsertValues(values, fileid, type_, keys)
		values = []string{}
		keys = []string{}
	}
	return count
}

// 递归处理文件
func dearDir(dirPath string) {
	dirPath = strings.Replace(dirPath, "/", "\\", -1)
	tmpBase64Path := "__tmp__.base64"
	eb := Console.EqBar{
		Size: 50,
	}
	fileChan := make(chan Node)
	dirChan := make(chan Node)
	next := make(chan int)
	over := make(chan int)
	fileId := GetUUID()
	DirDatabase.InsertDir(fileId, "", "", "dir", "", Tag)
	parentIdMap := map[string]string{
		dirPath: fileId,
	}
	dearNode := func(dir Node, nodeType string) {
		parentId := ""
		fileId := GetUUID()
		if dir.Parent != "" {
			parentId = parentIdMap[dir.Parent]
		}
		if nodeType == db.FileString {
			nodeType = path.Ext(dir.Name)
		}
		parentIdMap[dir.Parent+"\\"+dir.Name] = fileId
		if len(dir.Name) > 5 && strings.ToLower(dir.Name[len(dir.Name)-4:]) == db.DOCX {
			//fmt.Println("正在处理文件 -> " + dir.Name)
			/*total += */
			normalWord(path.Join(dir.Parent, dir.Name), tmpBase64Path)
			FileDatabase.InsertFile(fileId, DirAndFile.ReadAll(tmpBase64Path))
			nodeType = db.DOCX
		} else if len(dir.Name) > 5 && strings.ToLower(dir.Name[len(dir.Name)-4:]) == db.XLSX {
			//fmt.Println("正在处理文件 -> " + f.Name)
			/*total += */
			normalExcel(path.Join(dir.Parent, dir.Name), tmpBase64Path)
			FileDatabase.InsertFile(fileId, DirAndFile.ReadAll(tmpBase64Path))
			nodeType = db.XLSX
		}
		if dir.Name == dirPath {
			dir.Name = ""
		}
		DirDatabase.InsertDir(fileId, dir.Name, "", nodeType, parentId, "")
		next <- 1
	}
	go func() {
		for {
			select {
			case file := <-fileChan:
				dearNode(file, db.FileString)
			case dir := <-dirChan:
				dearNode(dir, db.DirString)
			case <-over:
				break
			}
		}
	}()
	fileNodes, dirNodes := anlize(dirPath)
	fmt.Println("开始写文件夹信息")
	eb.Total = len(dirNodes)
	eb.Percent = 0
	for _, dirNode := range dirNodes {
		dirChan <- dirNode
		eb.ShowNumPlusOne()
		<-next
	}
	fmt.Println("\r\n开始写文件信息")
	//fmt.Println("文件数量：",len(fileNodes))
	eb.Total = len(fileNodes)
	eb.Percent = 0
	for _, fileNode := range fileNodes {
		fileChan <- fileNode
		eb.ShowNumPlusOne()
		<-next
	}
	over <- 1
}

// 文件节点
type Node struct {
	Name   string
	Parent string
}

// 深度优先遍历文件夹
func anlize(path string) (fileNodes []Node, dirNodes []Node) {
	files := DirAndFile.GetSubDirOrFile(path)
	for _, file := range files {
		if !file.File {
			if file.Name == "System Volume Information" {
				continue
			}
			dirNodes = append(dirNodes, Node{
				Name:   file.Name,
				Parent: path,
			})
			_fileNodes, _dirNodes := anlize(path + "\\" + file.Name)
			fileNodes = append(fileNodes, _fileNodes...)
			dirNodes = append(dirNodes, _dirNodes...)
		} else {
			fileNodes = append(fileNodes, Node{
				Name:   file.Name,
				Parent: path,
			})
		}
	}
	return
}
