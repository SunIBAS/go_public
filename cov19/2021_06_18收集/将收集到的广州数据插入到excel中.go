package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"strconv"
	"strings"
)

var (
	xlsxFile *excelize.File
	startInd int
	//sheetName string = "数据变量"
	sheetName   string   = "sh"
	firstIndex  []string = []string{"B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L"}
	secondIndex []string = []string{"N", "O", "P", "Q", "R", "S", "T"}
)

func main() {
	dearXlsx()
	//nef := office.NewExcelFile{
	//	File: nil,
	//	Path: "D:\\20210413_xj\\tmp.xlsx",
	//}
	//nef.Init()
	//nef.File.NewSheet(sheetName)
	//xlsxFile = nef.File
	dear()
	xlsxFile.Save()
	//xlsxFile.
}

func dearXlsx() {
	// 读取文件，获取到应该从哪一行开始写
	f, err := excelize.OpenFile("D:\\20210413_xj\\数据情况广州all.xlsx")
	if err != nil {
		panic(err)
		return
	}
	xlsxFile = f
	startInd = 1
	for {
		v, _ := xlsxFile.GetCellValue(sheetName, "N"+strconv.Itoa(startInd))
		if v == "" {
			break
		}
		startInd++
	}
	fmt.Println(startInd)
}

func insertInToXlsxFirst(line string) {
	parts := strings.Split(line, "\t")
	lineInd := strconv.Itoa(startInd)
	for i := 0; i < len(parts); i++ {
		//
		//fmt.Println("axis = " + firstIndex[i] + lineInd)
		//fmt.Println("content = " + parts[i])
		xlsxFile.SetCellValue(sheetName, firstIndex[i]+lineInd, parts[i])
	}
}

func insertInToXlsxSecond(line string) {
	parts := strings.Split(line, "\t")
	lineInd := strconv.Itoa(startInd)
	for i := 0; i < len(parts); i++ {
		//
		//fmt.Println("axis = " + secondIndex[i] + lineInd)
		//fmt.Println("content = " + parts[i])
		xlsxFile.SetCellValue(sheetName, secondIndex[i]+lineInd, parts[i])
	}
	startInd++
}

func dear() {
	lines := DirAndFile.ReadAsFileAsLine("C:\\Users\\IBAS\\Downloads\\a.txt")
	ind := 0
	line := ""
	maxLen := len(lines) - 1
	for ind < maxLen {
		line = lines[ind]
		if len(line) > 2 && line[:2] == "==" {
			//进入新的行
			ind++ // 正常？
			if lines[ind] != "正常" {
				fmt.Println(lines[ind+1] + " >>>> 不正常")
			} else {
				fmt.Println(lines[ind+1] + " >>>> 正常")
			}
			ind++ // 首行
			insertInToXlsxFirst(lines[ind])
			ind++
			line = lines[ind]
			if len(line) == 0 {
				ind++
			} else {
				for len(line) > 2 && line[:2] != "==" && ind < maxLen {
					insertInToXlsxSecond(lines[ind])
					ind++
					line = lines[ind]
				}
			}
		}
	}
}
