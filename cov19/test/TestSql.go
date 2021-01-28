package main

import (
	"os"

	//"database/sql"
	//_ "github.com/mattn/go-sqlite3"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	sqlsint "public.sunibas.cn/go_public/cov19/sql"
	"strconv"
	"strings"
)

func main()  {
	//testAdding()
	//batchAddingAndWriteBackToDB(20210116,20210127)
	//saveToExcel()
	saveAllToExcel()
}

func testSelectDI()  {
	sql := sqlsint.SqlsInt{
		DbPath:   "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\cov19\\db\\db.db",
		Database: nil,
		ShowSql: true,
	}
	sql.Init()
	defer sql.Close()
	rows := sql.SelectDI("德国","20210115","daily_inj")
	jsonStr,_ := json.MarshalIndent(rows,"","")
	fmt.Print(string(jsonStr))
}

func testSelect()  {
	sql := sqlsint.SqlsInt{
		DbPath:   "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\cov19\\db\\db.db",
		Database: nil,
		ShowSql: true,
	}
	sql.Init()
	defer sql.Close()
	rows := sql.Select("date='20210114' and city_from='美国'","daily_inj")
	jsonStr,_ := json.MarshalIndent(rows,"","")
	fmt.Print(string(jsonStr))
}

func testAdding() {
	sql := sqlsint.SqlsInt{
		DbPath:   "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\cov19\\db\\db.db",
		Database: nil,
		ShowSql: true,
	}
	sql.Init()
	defer sql.Close()
	tbs := utils.Adding(sql,"20210114","20210115")
	jsonStr,_ := json.MarshalIndent(tbs,"","")
	fmt.Print(string(jsonStr))
}
// addingAndWriteBackToDB(from,end)
func batchAddingAndWriteBackToDB(fromDate,endDate int)  {
	sql := sqlsint.SqlsInt{
		DbPath:   "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\cov19\\db\\db.db",
		Database: nil,
		ShowSql: true,
	}
	sql.Init()
	defer sql.Close()
	//addingAndWriteBackToDB("20210126","20210127")
	for i := fromDate;i < endDate;i++ {
		pDate := strconv.Itoa(i)
		tDate := strconv.Itoa(i + 1)
		tbs := utils.Adding(sql,pDate,tDate)
		for _,d := range tbs {
			sql.Insert(d.CityFrom,d.CityTo,tDate,d.Num,"total_inj",d.Whos,d.Note)
		}
	}
}
func addingAndWriteBackToDB(pDate,tDate string)  {
	sql := sqlsint.SqlsInt{
		DbPath:   "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\cov19\\db\\db.db",
		Database: nil,
		ShowSql: true,
	}
	sql.Init()
	defer sql.Close()
	tbs := utils.Adding(sql,pDate,tDate)
	for _,d := range tbs {
		sql.Insert(d.CityFrom,d.CityTo,tDate,d.Num,"total_inj",d.Whos,d.Note)
	}
}

//
func saveToExcel() {
	sql := sqlsint.SqlsInt{
		DbPath:   "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\cov19\\db\\db.db",
		Database: nil,
		ShowSql: true,
	}
	sql.Init()
	defer sql.Close()
	utils.SaveToExcel(sql,"20210115",
		"D:\\Temp\\1204\\first.xlsx")
}
func TestDN()  {
	dn := Times.DaysNexter{
		Formater:"2006-01-02",
	}
	dn.Init(2020,1,30)
	fmt.Println(dn.NextDate())
	fmt.Println(dn.NextDate())
	fmt.Println(dn.NextDate())
	fmt.Println(dn.NextDate())
	fmt.Println(dn.NextDate())
}
func saveAllToExcel()  {
	saveExcelPath := "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\cov19\\_pre\\today.xlsx"
	os.Remove(saveExcelPath)
	sql := sqlsint.SqlsInt{
		DbPath:   "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\cov19\\db\\db.db",
		Database: nil,
		ShowSql: true,
	}
	sql.Init()
	defer sql.Close()
	f := excelize.NewFile()
	dn := Times.DaysNexter{
		Formater:"2006-01-02",
	}
	dn.Init(2020,9,29)
	for {
		if d,over := dn.NextDateToTD();
			over {
				break
		} else {
			d = strings.Replace(d,"-","",-1)
			utils.SaveToExcelF(sql,d,f)
		}
	}
	if err := f.SaveAs(saveExcelPath);err != nil {
		panic(err)
	}
	// 删除 sheet1
	ofile, err := excelize.OpenFile(saveExcelPath)
	if err != nil {
		panic(err)
	}
	ofile.DeleteSheet("Sheet1")
	if err := ofile.SaveAs(saveExcelPath); err != nil {
		println(err.Error())
	}
}