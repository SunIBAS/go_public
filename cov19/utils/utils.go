package utils

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	sqlsint "public.sunibas.cn/go_public/cov19/sql"
	"public.sunibas.cn/go_public/utils/Datas"
	"strconv"
)

// 计算累加
/**
 * si
 * pDate 前一天 20210114
 * tDate 当天   20210115
 */
func Adding(si sqlsint.SqlsInt,pDate,tDate string) []sqlsint.Table {
	pDatas := si.Select("date='" + pDate + "'","total_inj")
	tDatas := si.Select("date='" + tDate + "'","daily_inj")
	get := -1
	for _,d := range tDatas {
		get = -1
		for ind,dd := range pDatas {
			if d.CityFrom == dd.CityFrom && d.CityTo == dd.CityTo {
				get = ind
				break
			}
		}
		if get != -1 {
			pDatas[get].Num = strconv.Itoa(Datas.AtoI(pDatas[get].Num) + Datas.AtoI(d.Num))
		} else {
			pDatas = append(pDatas,d)
		}
	}
	for ind,_ := range pDatas {
		pDatas[ind].Note = "auto add"
		pDatas[ind].Whos = "auto add"
	}
	return pDatas
}

// 将数据保存到 excel 中，和之前保存的格式相同即可
// date 是保存哪一天的数据 20210115
func SaveToExcel(si sqlsint.SqlsInt,date,excelPath string) {
	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		panic(err)
	}
	SaveToExcelF(si,date,f)
}
func SaveToExcelF(si sqlsint.SqlsInt,date string,f * excelize.File) {
	tables := si.Select("date='" + date + "'","daily_inj")
	createSheetByTableOnly(date + "输入","当日(" + date + ")输入统计",f,tables)
	tables = si.Select("date='" + date + "'","total_inj")
	createSheetByTableOnly(date + "输入累计","累计输入(截至" + date + "的24时)",f,tables)
	tables = si.Select("date='" + date + "'","daily_wzz")
	createSheetByTableOnly(date + "无症状","当日(" + date + ")无症状统计",f,tables)
	tables = si.Select("date='" + date + "'","daily_ill")
	createSheetByTableOnly(date + "本土","当日(" + date + ")本土病例统计",f,tables)
	tables = si.Select("1=1 order by date ","daily_ill")
	dateInt := Datas.AtoI(date)
	for ind,t := range tables {
		if Datas.AtoI(t.Date) > dateInt {
			tables = tables[:ind]
			break
		}
	}
	createTotalIll(date + "本土累计","累计（截至" + date + "的24时)",f,tables)
	f.Save()
}
func createSheetByTableOnly(sheetName,title string,f * excelize.File,tbs []sqlsint.Table) {
	f.NewSheet(sheetName)
	f.SetCellValue(sheetName,"A1",title)
	f.SetCellValue(sheetName,"A2","省份")
	inMapper := map[string]int {} // city_to
	outMapper := map[string]int {} // city_from
	lastRow := 3
	curRow := 3
	lastCol := 2
	curCol := 2
	axis := ""
	for _,t := range tbs {
		if _,ok := inMapper[t.CityTo];ok {
			curRow = inMapper[t.CityTo]
		} else {
			curRow = lastRow
			lastRow = lastRow + 1
			inMapper[t.CityTo] = curRow
			axis,_ = excelize.CoordinatesToCellName(1,curRow)
			f.SetCellValue(sheetName,axis,t.CityTo)
		}
		if _,ok := outMapper[t.CityFrom];ok {
			curCol = outMapper[t.CityFrom]
		} else {
			curCol = lastCol
			lastCol = lastCol + 1
			outMapper[t.CityFrom] = curCol
			axis,_ = excelize.CoordinatesToCellName(curCol,2)
			f.SetCellValue(sheetName,axis,t.CityFrom)
		}
		axis,_ = excelize.CoordinatesToCellName(curCol,curRow)
		f.SetCellInt(sheetName,axis,Datas.AtoI(t.Num))
	}
	axis,_ = excelize.CoordinatesToCellName(lastCol,2)
	f.SetCellValue(sheetName,axis,"合计")
	axis,_ = excelize.CoordinatesToCellName(1,lastRow)
	f.SetCellValue(sheetName,axis,"合计")
	for row := 3;row < lastRow;row++ {
		axis,_ = excelize.CoordinatesToCellName(lastCol,row)
		faxis,_ := excelize.CoordinatesToCellName(2,row)
		taxis,_ := excelize.CoordinatesToCellName(lastCol - 1,row)
		f.SetCellFormula(sheetName,axis,"SUM(" + faxis + ":" + taxis + ")")
	}
	if lastCol != 2 {
		for col := 2;col <= lastCol;col++ {
			axis,_ = excelize.CoordinatesToCellName(col,lastRow)
			faxis,_ := excelize.CoordinatesToCellName(col,3)
			taxis,_ := excelize.CoordinatesToCellName(col,lastRow - 1)
			f.SetCellFormula(sheetName,axis,"SUM(" + faxis + ":" + taxis + ")")
		}
	}
	//excelize.CoordinatesToCellName()
}
func dateStrToDate(ds string) string {
	return ds[0:4] + "年" + ds[4:6] + "月" + ds[6:] + "日"
}
func createTotalIll(sheetName,title string,f * excelize.File,tbs []sqlsint.Table)  {
	f.NewSheet(sheetName)
	f.SetCellValue(sheetName,"B1",title)
	f.SetCellValue(sheetName,"A2","省份")
	inMapper := map[string]int {} // 年月日
	outMapper := map[string]int {} // 北京、广东
	lastRow := 3
	curRow := 3
	lastCol := 2
	curCol := 2
	axis := ""
	theDate := ""
	for _,t := range tbs {
		theDate = dateStrToDate(t.Date)
		if _,ok := inMapper[theDate];ok {
			curRow = inMapper[theDate]
		} else {
			curRow = lastRow
			lastRow = lastRow + 1
			inMapper[theDate] = curRow
			axis,_ = excelize.CoordinatesToCellName(1,curRow)
			f.SetCellValue(sheetName,axis,theDate)
		}
		if _,ok := outMapper[t.CityTo];ok {
			curCol = outMapper[t.CityTo]
		} else {
			curCol = lastCol
			lastCol = lastCol + 1
			outMapper[t.CityTo] = curCol
			axis,_ = excelize.CoordinatesToCellName(curCol,2)
			f.SetCellValue(sheetName,axis,t.CityTo)
		}
		axis,_ = excelize.CoordinatesToCellName(curCol,curRow)
		f.SetCellInt(sheetName,axis,Datas.AtoI(t.Num))
	}
	axis,_ = excelize.CoordinatesToCellName(lastCol,2)
	f.SetCellValue(sheetName,axis,"合计")
	axis,_ = excelize.CoordinatesToCellName(1,lastRow)
	f.SetCellValue(sheetName,axis,"合计")
	for row := 3;row < lastRow;row++ {
		axis,_ = excelize.CoordinatesToCellName(lastCol,row)
		faxis,_ := excelize.CoordinatesToCellName(2,row)
		taxis,_ := excelize.CoordinatesToCellName(lastCol - 1,row)
		f.SetCellFormula(sheetName,axis,"SUM(" + faxis + ":" + taxis + ")")
	}
	if 2 != lastCol {
		for col := 2;col <= lastCol;col++ {
			axis,_ = excelize.CoordinatesToCellName(col,lastRow)
			faxis,_ := excelize.CoordinatesToCellName(col,3)
			taxis,_ := excelize.CoordinatesToCellName(col,lastRow - 1)
			f.SetCellFormula(sheetName,axis,"SUM(" + faxis + ":" + taxis + ")")
		}
	}
}