package main

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"log"
	"os"
	"regexp"
	"strings"
)

var logFile, _ = os.OpenFile("C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\cov19\\_pre\\log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
func main()  {
	log.SetOutput(logFile)
	readExcel()
}

// str = "1.14本土累计"
// tar = "本土" "本土累计"
// length = 2 4
// return false true
func checkEndsWith(str,tar string,length int) bool {
	return len(str) > 3 * length && str[len(str) - 3 * length:] == tar
}

func readExcel()  {
	f, err := excelize.OpenFile("C:\\Users\\IBAS\\Desktop\\bak.xlsx")
	if err != nil {
		panic(err)
		return
	}
	for _, name := range f.GetSheetMap() {
		//if len(name) > 6 && name[len(name) - 6:] == "输入" {
		if checkEndsWith(name,"输入",2) {
			class_1(f,name,"daily_inj")
		} else if len(name) > 12 && name[len(name) - 12:] == "输入累计" {
			class_1(f,name,"total_inj")
		} else if len(name) > 9 && name[len(name) - 9:] == "无症状" {
			class_1(f,name,"daily_wzz")
		} else if checkEndsWith(name,"累计本土",4) {
			//class_1(f,name,"daily_ill")
			continue
		} else if checkEndsWith(name,"本土",2) {
			class_1(f,name,"daily_ill")
		}
	}
}
var numReg,_ = regexp.Compile("^[0-9]+$")
func class_1(f * excelize.File,sheetName string,_type string) {
	log.Println("--start--")
	rows,_ := f.GetRows(sheetName)
	start := -1
	end := 0
	line := []string{}
	for ri, row := range rows {
		if start != -1 {} else if row[0] == "省份" {
			start = ri
			for ind,col := range row {
				if col == "合计" || col == "小计" {
					end = ind
					break
				}
			}
			continue
		} else {
			log.Println(_type)
			log.Println(strings.Join(row," "))
			continue
		}
		i := 1
		if row[0] == "合计" || row[0] == "小计" {
			break
		}
		for ;i < end;i++ {
			if numReg.MatchString(row[i]) {
				line = []string{}
				line = append(line, rows[start][i])
				line = append(line, row[0])
				line = append(line, row[i])
				log.Println(strings.Join(line,"\t"))
			}
		}
	}
	log.Println("--end--")
}
