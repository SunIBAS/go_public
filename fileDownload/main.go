package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"public.sunibas.cn/go_public/utils/DirAndFile"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("\t fileDownload.exe url file [t]")
		fmt.Println("\t url\t是文件的链接")
		fmt.Println("\t file\t是文件的完整保存路径")
		fmt.Println("\t t\t表示如果存在则覆盖，否则存在则跳过")
		fmt.Println("\t ")
		return
	}
	imgUrl := os.Args[1]
	// http://www.cawater-info.net/pdf/pyanj-river-morphology.pdf
	// http://www.cawater-info.net/pdf/pyanj-river-morphology.pdf
	savePath := os.Args[2] // "D:\\Temp\\cawater-info\\tmp\\pdf\\"
	recoverFile := false
	if len(os.Args) == 4 && os.Args[3] == "t" {
		recoverFile = true
	}
	if !DirAndFile.PathNotExists(savePath) {
		if !recoverFile {
			fmt.Println("文件 " + savePath + " 已经存在")
			return
		}
	}

	dir := filepath.Dir(savePath)

	// Get the data
	resp, err := http.Get(imgUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if _, e := os.Stat(dir); e != nil {
		os.Mkdir(dir, os.ModePerm)
	}

	// 创建一个文件用于保存
	out, err := os.Create(savePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}
