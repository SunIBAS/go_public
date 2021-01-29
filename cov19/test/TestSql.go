package main

import (
	"os"
	"public.sunibas.cn/go_public/cov19/utils"
	"public.sunibas.cn/go_public/utils/Times"
	"public.sunibas.cn/go_public/utils/office"

	//"database/sql"
	//_ "github.com/mattn/go-sqlite3"
	"encoding/json"
	"fmt"
	sqlsint "public.sunibas.cn/go_public/cov19/sql"
	"strconv"
	"strings"
)

func main() {
	//testAdding()
	//batchAddingAndWriteBackToDB(20210116,20210127)
	//saveToExcel()
	saveAllToExcel()
}

func testSelectDI() {
	sql := sqlsint.SqlsInt{
		DbPath:   "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\cov19\\db\\db.db",
		Database: nil,
		ShowSql:  true,
	}
	sql.Init()
	defer sql.Close()
	rows := sql.SelectDI("德国", "20210115", "daily_inj")
	jsonStr, _ := json.MarshalIndent(rows, "", "")
	fmt.Print(string(jsonStr))
}

func testSelect() {
	sql := sqlsint.SqlsInt{
		DbPath:   "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\cov19\\db\\db.db",
		Database: nil,
		ShowSql:  true,
	}
	sql.Init()
	defer sql.Close()
	rows := sql.Select("date='20210114' and city_from='美国'", "daily_inj")
	jsonStr, _ := json.MarshalIndent(rows, "", "")
	fmt.Print(string(jsonStr))
}

func testAdding() {
	sql := sqlsint.SqlsInt{
		DbPath:   "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\cov19\\db\\db.db",
		Database: nil,
		ShowSql:  true,
	}
	sql.Init()
	defer sql.Close()
	tbs := utils.Adding(sql, "20210114", "20210115")
	jsonStr, _ := json.MarshalIndent(tbs, "", "")
	fmt.Print(string(jsonStr))
}

// addingAndWriteBackToDB(from,end)
func batchAddingAndWriteBackToDB(fromDate, endDate int) {
	sql := sqlsint.SqlsInt{
		DbPath:   "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\cov19\\db\\db.db",
		Database: nil,
		ShowSql:  true,
	}
	sql.Init()
	defer sql.Close()
	//addingAndWriteBackToDB("20210126","20210127")
	for i := fromDate; i < endDate; i++ {
		pDate := strconv.Itoa(i)
		tDate := strconv.Itoa(i + 1)
		tbs := utils.Adding(sql, pDate, tDate)
		for _, d := range tbs {
			sql.Insert(d.CityFrom, d.CityTo, tDate, d.Num, "total_inj", d.Whos, d.Note)
		}
	}
}
func addingAndWriteBackToDB(pDate, tDate string) {
	sql := sqlsint.SqlsInt{
		DbPath:   "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\cov19\\db\\db.db",
		Database: nil,
		ShowSql:  true,
	}
	sql.Init()
	defer sql.Close()
	tbs := utils.Adding(sql, pDate, tDate)
	for _, d := range tbs {
		sql.Insert(d.CityFrom, d.CityTo, tDate, d.Num, "total_inj", d.Whos, d.Note)
	}
}

//
func saveToExcel() {
	sql := sqlsint.SqlsInt{
		DbPath:   "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\cov19\\db\\db.db",
		Database: nil,
		ShowSql:  true,
	}
	sql.Init()
	defer sql.Close()
	utils.SaveToExcel(sql, "20210115",
		"D:\\Temp\\1204\\first.xlsx")
}
func TestDN() {
	dn := Times.DaysNexter{
		Formater: "2006-01-02",
	}
	dn.Init(2020, 1, 30)
	fmt.Println(dn.NextDate())
	fmt.Println(dn.NextDate())
	fmt.Println(dn.NextDate())
	fmt.Println(dn.NextDate())
	fmt.Println(dn.NextDate())
}
func saveAllToExcel() {
	saveExcelPath := "C:\\Users\\IBAS\\go\\src\\public.sunibas.cn\\go_public\\cov19\\_pre\\today.xlsx"
	os.Remove(saveExcelPath)
	sql := sqlsint.SqlsInt{
		DbPath:   "C:\\Users\\IBAS\\go\\src\\public.sunibas.cn\\go_public\\cov19\\db\\db.db",
		Database: nil,
		ShowSql:  true,
	}
	sql.Init()
	defer sql.Close()
	nef := office.NewExcelFile{
		File: nil,
		Path: saveExcelPath,
	}
	nef.Init()
	dn := Times.DaysNexter{
		Formater: "2006-01-02",
	}
	dn.Init(2020, 9, 29)
	for {
		if d, over := dn.NextDateToTD(); over {
			break
		} else {
			d = strings.Replace(d, "-", "", -1)
			utils.SaveToExcelF(sql, d, nef.File)
		}
	}
	nef.SaveAndDeleteSheet1()
}
