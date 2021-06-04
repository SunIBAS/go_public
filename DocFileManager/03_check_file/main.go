package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"public.sunibas.cn/go_public/DocFileManager/utilite"
	"public.sunibas.cn/go_public/utils/DirAndFile"
)

func main() {
	config := utilite.GetConfig()
	_database1, _err1 := sql.Open("sqlite3", config.DirDbPath)
	defer _database1.Close()
	if _err1 != nil {
		panic(_err1)
	}
	check(_database1)
}

func check(db *sql.DB) {
	dirs := utilite.Sql_select_from_dir(db, " where status='"+utilite.FileStatus_add+"' and type='file';")
	for _, d := range dirs {
		oType := "." + DirAndFile.CheckType()
	}
}
