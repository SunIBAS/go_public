package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"public.sunibas.cn/go_public/translateOffice/InsertJson"
	"public.sunibas.cn/go_public/translateOffice/NormalFile"
	"public.sunibas.cn/go_public/translateOffice/Server"
	"public.sunibas.cn/go_public/translateOffice/WriteOutTraned"
	"public.sunibas.cn/go_public/translateOffice/db"
)

var programName = "translateOffice"

func usage() {
	fmt.Println("1.格式化并生成文档模板")
	fmt.Println("\t将 inDir 下的所有文档格式化后写到 outDir 中")
	fmt.Println("\t如果 outDir 为 [db]，表示将文件暂存到 【数据库】 文件中")
	fmt.Println("\t[db] 模式下将递归遍历 inDir 目录")
	fmt.Println("\t非 [db] 模式不递归")
	fmt.Println("(不递归)>>> " + programName + " -format inDir outDir")
	fmt.Println("(递归)>>> " + programName + " -format inDir [db] tag")
	fmt.Println()
	fmt.Println("2.翻译，将开启一个服务，用于在浏览器上翻译文本")
	fmt.Println(">>>" + programName + " -translate 8080")
	fmt.Println()
	fmt.Println("3.写回第二步生成的 json 文件")
	fmt.Println(">>>" + programName + " -insert json")
	fmt.Println()
	fmt.Println("4.写出翻译好的文件")
	fmt.Println("\t需要提供第一步的 outDir 作为 base64Dir")
	fmt.Println("\t这里的 outDir 是翻译好的文档的输出文件夹")
	fmt.Println("\t如果第一步 outDir 为 [db] 这里也使用 [db]")
	fmt.Println("\t使用 [db] 将会把全部文件写出，并创建递归目录")
	fmt.Println(">>>" + programName + " -ok base64Dir outDir")
	fmt.Println(">>>" + programName + " -ok [db] outDir tag")
	fmt.Println("\t[db] 条件下使用 :test 可以测试写出文件目录是否正常")
	fmt.Println(">>>" + programName + " -ok [db] outDir tag:test")
}

// -format C:\Users\IBAS\Desktop\tran\docx [db] tag
// -translate 8080
// -insert D:\Temp\doc\base64\1.json
// -ok D:\Temp\doc\base64\ D:\Temp\doc\trans\
// -ok [db] D:\Temp\doc\trans\out test1
func main() {
	langDb := db.LangDb{}
	langDb.Init()
	dirDb := db.DirDb{}
	dirDb.Init()
	fileDb := db.FileDb{}
	fileDb.Init()
	defer langDb.LangDb.Close()
	defer dirDb.DirDb.Close()
	defer fileDb.FileDb.Close()

	if len(os.Args) > 2 {
		if os.Args[1] == "-format" {
			if os.Args[3] == "[db]" {
				if len(os.Args) == 4 {
					NormalFile.NormalFile(&langDb, &dirDb, &fileDb, os.Args[2], os.Args[3], NormalFile.GetUUID())
					return
				} else if len(os.Args) == 5 {
					NormalFile.NormalFile(&langDb, &dirDb, &fileDb, os.Args[2], os.Args[3], os.Args[4])
					return
				}
			} else {
				if len(os.Args) == 4 {
					NormalFile.NormalFile(&langDb, nil, nil, os.Args[2], os.Args[3], NormalFile.GetUUID())
					return
				} else if len(os.Args) == 5 {
					NormalFile.NormalFile(&langDb, nil, nil, os.Args[2], os.Args[3], os.Args[4])
					return
				}
			}
		} else if os.Args[1] == "-translate" {
			if len(os.Args) == 3 {
				Server.StartServer(&langDb, os.Args[2])
				return
			}
		} else if os.Args[1] == "-insert" {
			if len(os.Args) == 3 {
				InsertJson.InsertJson(&langDb, os.Args[2])
				return
			}
		} else if os.Args[1] == "-ok" {
			if len(os.Args) == 4 {
				WriteOutTraned.WriteOut(&langDb, &dirDb, &fileDb, os.Args[2], os.Args[3], "")
				return
			} else if len(os.Args) == 5 {
				WriteOutTraned.WriteOut(&langDb, &dirDb, &fileDb, "[db]", os.Args[3], os.Args[4])
				return
			}
		}
	}
	usage()
}
