package WriteOutTraned

import (
	"fmt"
	"os"
	"path"
	"public.sunibas.cn/go_public/translateOffice/db"
	"public.sunibas.cn/go_public/utils/Console"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"public.sunibas.cn/go_public/utils/ParseOffice"
)

var tmpFilePath = "__tmp__file__.tmp"
var idLength = 36
var langDatabase *db.LangDb
var dirDatabase *db.DirDb
var fileDatabase *db.FileDb
var TestOut = false

func WriteOut(
	_lang_db_ *db.LangDb,
	_dir_db_ *db.DirDb,
	_file_db_ *db.FileDb, base64Dir, outDir, tag string) {
	langDatabase = _lang_db_
	dirDatabase = _dir_db_
	fileDatabase = _file_db_
	if base64Dir != "[db]" {
		files := DirAndFile.GetSubDirOrFile(base64Dir)
		eb := Console.EqBar{
			Total:   len(files),
			Percent: 0,
			Size:    50,
		}
		for _, f := range files {
			if f.File {
				if len(f.Name) > 7 && f.Name[len(f.Name)-6:] == "base64" {
					base64 := DirAndFile.ReadAll(path.Join(base64Dir, f.Name))
					if len(base64) > len(db.DOCX) && base64[0:len(db.DOCX)] == db.DOCX {
						fmt.Println("正在处理文件 -> " + f.Name)
						RecoverWord(base64[4:], path.Join(outDir, f.Name[0:len(f.Name)-6]+db.DOCX))
					} else if len(base64) > len(db.XLSX) && base64[0:len(db.XLSX)] == db.XLSX {
						fmt.Println("正在处理文件 -> " + f.Name)
						RecoverExcel(base64[4:], path.Join(outDir, f.Name[0:len(f.Name)-6]+db.XLSX))
					}
				}
				eb.ShowNumPlusOne()
			}
		}
	} else {
		fmt.Println("tag = " + tag)
		root := dirDatabase.SelectByTagFromDir(tag)
		if len(tag) > len(":test") && tag[len(tag)-len(":test"):] == ":test" {
			fmt.Println("test")
			tag = tag[:len(tag)-len(":test")]
			root = dirDatabase.SelectByTagFromDir(tag)
			TestOut = true
		}
		if len(root) == 1 {
			recursionDB(root[0], outDir)
		}
	}
}

func recursionDB(root db.Dir, rootPath string) {
	if root.Type == db.DOCX {
		files := fileDatabase.SelectByIdFromFile(root.Id)
		if len(files) == 1 {
			fmt.Println("写出 word 文件 " + root.Name)
			if TestOut {
				DirAndFile.WriteWithWriteString(path.Join(rootPath, root.Name), root.Id)
			} else {
				RecoverWord(files[0].Content[4:], path.Join(rootPath, root.Name))
			}
		}
	} else if root.Type == db.XLSX {
		files := fileDatabase.SelectByIdFromFile(root.Id)
		if len(files) == 1 {
			fmt.Println("写出 excel 文件 " + root.Name)
			if TestOut {
				DirAndFile.WriteWithWriteString(path.Join(rootPath, root.Name), root.Id)
			} else {
				RecoverExcel(files[0].Content[4:], path.Join(rootPath, root.Name))
			}
		}
	} else {
		if len(root.Name) != 0 {
			rootPath = path.Join(rootPath, root.Name)
		}
		if t, _ := DirAndFile.PathExistsAndType(rootPath); t != 1 {
			fmt.Println("创建文件夹 " + rootPath)
			os.Mkdir(rootPath, os.ModeDir)
		}
		dirs := dirDatabase.SelectByParentIdFromDir(root.Id)
		for _, dir := range dirs {
			recursionDB(dir, rootPath)
		}
	}
}

func RecoverWord(base64, outPath string) {
	fileId := base64[0:36]
	fileContent := base64[36:]
	DirAndFile.SaveBase64StringToByteFile(fileContent, tmpFilePath)
	langs := langDatabase.Select(fileId, "")
	rp := ParseOffice.RpWord{
		InFile:  tmpFilePath,
		OutFile: outPath,
		XmlName: "",
		Words:   nil,
		Xml:     "",
		Merge:   false,
	}
	rp.Init()
	repWords := map[string]string{}
	for _, lang := range langs {
		repWords[`{{`+lang.Id+`}}`] = lang.Tran //`第` + strconv.Itoa(ind) + `条`
	}
	rp.Replace(repWords)
	rp.WriteOut()
}
func RecoverExcel(base64, outPath string) {
	fileId := base64[0:36]
	fileContent := base64[36:]
	DirAndFile.SaveBase64StringToByteFile(fileContent, tmpFilePath)
	langs := langDatabase.Select(fileId, "")
	rp := ParseOffice.RpExcel{
		InFile:  tmpFilePath,
		OutFile: outPath,
		XmlName: "",
		Words:   nil,
		Xml:     "",
		Merge:   false,
	}
	rp.Init()
	repWords := map[string]string{}
	for _, lang := range langs {
		if lang.Type == db.WorkBook {
			repWords[`{{`+lang.Id+`}}`] = lang.Tran //`第` + strconv.Itoa(ind) + `条`
		}
	}
	rp.Replace(repWords)
	repWords = map[string]string{}
	for _, lang := range langs {
		if lang.Type == db.SheetName {
			repWords[lang.Id] = lang.Tran //`第` + strconv.Itoa(ind) + `条`
		}
	}
	rp.ReplaceSheetNameFromUUID(repWords)

	rp.WriteOut()
}
