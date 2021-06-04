package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
	"public.sunibas.cn/go_public/DocFileManager/utilite"
	"public.sunibas.cn/go_public/utils/Console"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"public.sunibas.cn/go_public/utils/SqliteSql"
	"strings"
)

var id, dirpath, name, token string

/**
 * 需要参数有
 * 1. id 默认为 0 表示挂载到顶级目录下（在逻辑上的父级目录的id）
 * 2. dirpath 要加入的文件夹 绝对路径
 * 3. name 逻辑上的文件夹名称（在逻辑上显示当前加入文件夹的文件夹名称）
 * 		如果只需要将文件加入到指定目录下，则 name 使用 -
 * 4. token 唯一标识，可以用于删除该批次下的所有文件（用于数据回滚，或者错误操作、断电导致的错误）
 * 		如果不知道写什么可以用 -
 */
func main() {
	if utilite.CheckConfigFileExist() {
		if len(os.Args) == 5 {
			id = os.Args[1]
			dirpath = os.Args[2]
			name = os.Args[3]
			token = os.Args[4]
			addDirs()
		} else {
			fmt.Println("需要三个参数")
		}
	} else {
		fmt.Println("配置文件不存在，请先进行初始化")
	}
}

func addDirs() {
	config := utilite.GetConfig()
	_database1, _err1 := sql.Open("sqlite3", config.DirDbPath)
	defer _database1.Close()
	if _err1 != nil {
		panic(_err1)
	}
	if token == "-" {
		token = utilite.CreateId()
	}

	initDB(dirpath, _database1, func(s string) {}, id, token)
}

// 遍历生成整个文件夹的层级数据库
// 将文件重命名并复制到指定位置
// 生成相应的数据库语句
// dirPath 要处理的文件夹
// defaultParentId 表示将文件写入到指定的目录下(数据库中的文件id)
// note 标记，用于对新加入的文件做标记便于后期处理
func initDB(dirPath string, database *sql.DB, returnInformation func(string), defaultParentId, note string) {
	eb := Console.EqBar{
		Size: 50,
	}
	fileChan := make(chan Node)
	dirChan := make(chan Node)
	next := make(chan int)
	over := make(chan int)
	dirPath = strings.Replace(dirPath, "/", "\\", -1)
	parentIdMap := map[string]string{
		dirPath + "\\": defaultParentId,
	}
	langIdMap := map[string]string{}
	if dirPath[len(dirPath)-1] != '\\' {
		dirPath += "\\"
	}
	//dearFile := func(file Node) {
	//	//
	//	next <- 1
	//}
	dearNode := func(dir Node, nodeType string) {
		// 先判断父文件夹的情况
		parentId := ""
		langId := ""
		if dir.Parent != "" {
			parentId = parentIdMap[dir.Parent]
		}
		// 语言包 id
		if _, ok := langIdMap[dir.Name]; ok {
			langId = langIdMap[dir.Name]
		} else {
			langId = utilite.CreateId()
			SqliteSql.ExecSqlString(database, utilite.Sql_insert_into_lang(langId, utilite.OrgType, dir.Name))
			langIdMap[dir.Name] = langId
		}

		dirId := utilite.CreateId()
		if nodeType == "file" {
			SqliteSql.ExecSqlString(database, utilite.Sql_insert_into_dir(dirId, parentId, dir.Name, nodeType, "0", langId, strings.ToLower(path.Ext(dir.Name)), dir.Parent, note, utilite.FileStatus_add))
		} else {
			SqliteSql.ExecSqlString(database, utilite.Sql_insert_into_dir(dirId, parentId, dir.Name, nodeType, "0", langId, "", dir.Parent, note, utilite.FileStatus_add))
			parentIdMap[dir.Parent+"\\"+dir.Name] = dirId
		}
		next <- 1
	}
	go func() {
		for {
			select {
			case file := <-fileChan:
				dearNode(file, "file")
			case dir := <-dirChan:
				dearNode(dir, "dir")
			case <-over:
				returnInformation("over")
				break
			}
		}
	}()
	fileNodes, dirNodes := anlize(dirPath)
	fmt.Println("开始写文件夹信息")
	returnInformation("开始写文件夹信息")
	eb.Total = len(dirNodes)
	eb.Percent = 0
	for _, dirNode := range dirNodes {
		dirChan <- dirNode
		eb.ShowNumPlusOne()
		returnInformation(eb.GetShowNum())
		<-next
	}
	fmt.Println("\r\n开始写文件信息")
	returnInformation("开始写文件信息")
	//fmt.Println("文件数量：",len(fileNodes))
	eb.Total = len(fileNodes)
	eb.Percent = 0
	for _, fileNode := range fileNodes {
		fileChan <- fileNode
		eb.ShowNumPlusOne()
		returnInformation(eb.GetShowNum())
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
