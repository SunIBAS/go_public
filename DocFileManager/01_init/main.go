package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"public.sunibas.cn/go_public/DocFileManager/utilite"
	"public.sunibas.cn/go_public/utils/Console"
	"public.sunibas.cn/go_public/utils/DirAndFile"
)

func main() {
	var input string
	config := utilite.Config{
		DirDbPath:  "",
		LangDbPath: "",
		TmpPath:    "",
		Port:       "",
	}
	// 检查配置是否存在
	if utilite.CheckConfigFileExist() {
		Console.InputLine(&input, "y", "配置文件已存在，是否重新配置(y/n)")
		if input == "y" {
			os.Remove(utilite.ConfigFileName)
		} else {
			return
		}
	}
	Console.InputLine(&input, "dir.db", "数据库1，用于保存文件结构")
	config.DirDbPath = input
	Console.InputLine(&input, "lang.db", "数据库2，用于保存文件模板和翻译内容")
	config.LangDbPath = input
	Console.InputLine(&input, "C:\\Users\\IBAS\\Documents\\tmp", "临时目录，用于保存处理过程中的临时文件")
	config.TmpPath = input
	Console.InputLine(&input, "8080", "启动服务端口")
	config.Port = input

	_database1, _err1 := sql.Open("sqlite3", config.DirDbPath)
	defer _database1.Close()
	if _err1 != nil {
		panic(_err1)
	}
	utilite.CreateDirTable(_database1)

	_database2, _err2 := sql.Open("sqlite3", config.LangDbPath)
	defer _database2.Close()
	if _err2 != nil {
		panic(_err2)
	}
	utilite.CreateLangTable(_database2)

	if content, err := json.Marshal(config); err == nil {
		DirAndFile.WriteWithIOUtilByte(utilite.ConfigFileName, content)
	} else {
		panic(err)
	}
}
