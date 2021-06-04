package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
	"path/filepath"
	"public.sunibas.cn/go_public/dearConv19/db"
	"public.sunibas.cn/go_public/utils/Datas"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"regexp"
	"strconv"
	"strings"
)

var (
	database db.Db
)

func help() {
	fmt.Println("插入一个文件")
	fmt.Println("dearConv19 -if d:\\a\\01-02-2021.csv")
	fmt.Println()
	fmt.Println("插入一个文件夹下的所有文件")
	fmt.Println("dearConv19 -id d:\\a")
	fmt.Println()
	fmt.Println("写出某个国家（地区）的数据，province 如果为空表示该国家的总计 province=__total__")
	fmt.Println("dearConv19 -o[flag] outDir country [province]")
	fmt.Println("dearConv19 -ocdr d:\\a country province")
	fmt.Println("生成的文件为 d:\\a\\country_province.csv")
	fmt.Println("flag = cdrCDR 六个字符表示")
	fmt.Println("\tc 确诊\td 死亡\tr 治愈")
	fmt.Println("\tC 确诊差值\tD 死亡差值\tR 治愈差值")
}

func main() {
	database = db.Db{}
	database.Init()
	defer database.Db.Close()
	ok := false

	if len(os.Args) > 1 {
		if os.Args[1] == "-if" {
			insertOneCsv(os.Args[2])
			ok = true
		} else if os.Args[1] == "-id" {
			insertFromDirCsv(os.Args[2])
			ok = true
		} else if os.Args[1][0:2] == "-o" {
			fmt.Println(os.Args[1][2:])
			if len(os.Args) == 4 {
				writeOut(os.Args[2], os.Args[3], "__total__", os.Args[1][2:])
				ok = true
			} else if len(os.Args) == 5 {
				writeOut(os.Args[2], os.Args[3], os.Args[4], os.Args[1][2:])
				ok = true
			}
		}
	}
	if !ok {
		help()
	}

	//csvPath := "D:\\20210413_xj\\COVID-19-master\\csse_covid_19_data\\csse_covid_19_daily_reports\\01-22-2020.csv"
	//insertOneCsv(csvPath)
	//csvPath = "D:\\20210413_xj\\COVID-19-master\\csse_covid_19_data\\csse_covid_19_daily_reports\\02-01-2020.csv"
	//insertOneCsv(csvPath)
	//csvPath := "D:\\20210413_xj\\COVID-19-master\\csse_covid_19_data\\csse_covid_19_daily_reports"
	//insertFromDirCsv(csvPath)

	//savePath := "D:\\20210413_xj\\COVID-19-master\\get"
	//countries := []string {"Cuba", "Bolivia", "Afghanistan", "Argentina", "Austria", "New Zealand", "Armenia", "India", "Peru", "Algeria", "Australia", "Bangladesh", "Belarus", "Philippines", "Azerbaijan", "Chile", "Belgium", "Bosnia and Herzegovina", "Singapore", "Bahrain", "South Africa", "France", "Bulgaria", "Thailand", "China", "Germany", "Canada", "Malaysia", "Djibouti", "Greece", "Croatia", "Panama", "Egypt", "Guinea", "Czech Republic", "Nigeria", "Iran", "Indonesia", "Denmark", "Puerto Rico", "Iraq", "Ireland", "Estonia", "Cameroon", "Israel", "Italy", "Finland", "Colombia", "Kazakhstan", "Japan", "Hungary", "Dominican Republic", "Kuwait", "Luxembourg", "Lithuania", "Ecuador", "Mexico", "Netherlands", "Moldova", "Brazil", "Morocco", "Portugal", "Norway", "Ghana", "Oman", "United Kingdom", "Poland", "Pakistan", "Romania", "Qatar", "Russia", "Saudi Arabia", "Serbia", "Spain", "Slovakia", "Turkey", "Slovenia", "United Arab Emirates", "South Korea", "Iceland", "US", "Sweden", "Uzbekistan", "Switzerland", "Ukraine"}
	////countries := []string {"US"}
	//for _,country := range countries {
	//	writeOut(savePath,country,"__total__")
	//}
}

func writeOut(savepath, country, province, flag string) {
	header := []string{"Confirmed", "Deaths", "Recovered", "Confirmed_deta", "Deaths_deta", "Recovered_deta"}
	headerLine := "Date"
	lines := []string{}
	last := []int{-1, -1, -1}
	dataIndexMap := map[string]int{
		"c": 0,
		"d": 1,
		"r": 2,
		"C": 3,
		"D": 4,
		"R": 5,
	}
	dataIndex := []int{}
	for i := 0; i < len(flag); i++ {
		dataIndex = append(dataIndex, dataIndexMap[string(flag[i])])
		headerLine += "," + header[dataIndexMap[string(flag[i])]]
	}
	lines = append(lines, headerLine)
	fmt.Println("write out " + country + province)
	records := database.Select(country, province)
	total := len(records)
	for ind := 0; ind < total; ind++ {
		curLineData := []string{}
		curLine := records[ind].Date
		confirmed := Datas.AtoI(records[ind].Confirmed)
		deaths := Datas.AtoI(records[ind].Deaths)
		recovered := Datas.AtoI(records[ind].Recovered)
		for {
			if ind+1 < total && records[ind].Date == records[ind+1].Date {
				ind++
				confirmed += Datas.AtoI(records[ind].Confirmed)
				deaths += Datas.AtoI(records[ind].Deaths)
				recovered += Datas.AtoI(records[ind].Recovered)
			} else {
				break
			}
		}
		if last[0] == -1 {
			curLineData = []string{
				strconv.Itoa(confirmed),
				strconv.Itoa(deaths),
				strconv.Itoa(recovered),
				"0", "0", "0",
			}
		} else {
			curLineData = []string{
				strconv.Itoa(confirmed),
				strconv.Itoa(deaths),
				strconv.Itoa(recovered),
				strconv.Itoa(confirmed - last[0]),
				strconv.Itoa(deaths - last[1]),
				strconv.Itoa(recovered - last[2]),
			}
		}
		last[0] = confirmed
		last[1] = deaths
		last[2] = recovered
		for i := 0; i < len(dataIndex); i++ {
			curLine += "," + curLineData[dataIndex[i]]
		}
		lines = append(lines, curLine)
	}
	var saveFullPath string
	if province == "__total__" {
		saveFullPath = path.Join(savepath, country+".csv")
	} else {
		saveFullPath = path.Join(savepath, country+"."+province+".csv")
	}
	DirAndFile.WriteWithWriteString(saveFullPath, strings.Join(lines, "\r\n"))
}

/**
 * 将某个文件夹下的所有 csv 文件插入到数据库中
 */
func insertFromDirCsv(csvPath string) {
	files := DirAndFile.GetSubDirOrFile(csvPath)
	var csvPathReg, _ = regexp.Compile("^[0-9]{2}-[0-9]{2}-[0-9]{4}.csv$")
	for _, f := range files {
		if csvPathReg.MatchString(f.Name) {
			insertOneCsv(path.Join(csvPath, f.Name))
		}
	}
}

func splitLine(line string, zeros []int) []string {
	line = strings.Replace(line, "'", "''", -1)
	ws := strings.Split(line, ",")
	if strings.Contains(line, "\"") {
		// "Chicago, IL",US,2/1/2020 19:43,2,0,0
		parts := strings.Split(line, "\"")
		ws = []string{}
		for ind, part := range parts {
			if ind%2 == 1 {
				ws = append(ws, part)
			} else {
				if len(part) == 0 {
					continue
				}
				if ind == 0 {
					ws = append(ws, strings.Split(part[0:len(part)-1], ",")...)
				} else if ind == len(parts)-1 {
					ws = append(ws, strings.Split(part[1:], ",")...)
				} else {
					ws = append(ws, strings.Split(part[1:len(part)-1], ",")...)
				}
			}
		}
	}
	for _, z := range zeros {
		if len(ws[z]) == 0 {
			ws[z] = "0"
		}
	}
	return ws
}
func insertOneCsv(csvPath string) {
	// month-day-year.csv'
	filepath.Base(csvPath)
	fns := strings.Split(strings.Split(filepath.Base(csvPath), ".")[0], "-")
	lines := DirAndFile.ReadAsFileAsLine(csvPath)
	records := map[string]int{}
	recordsCountry := []string{}
	Date := fns[2] + fns[0] + fns[1]
	lines[0] = strings.TrimSpace(lines[0])
	if lines[0] == "FIPS,Admin2,Province_State,Country_Region,Last_Update,Lat,Long_,Confirmed,Deaths,Recovered,Active,Combined_Key,Incident_Rate,Case_Fatality_Ratio" ||
		lines[0] == "FIPS,Admin2,Province_State,Country_Region,Last_Update,Lat,Long_,Confirmed,Deaths,Recovered,Active,Combined_Key,Incidence_Rate,Case-Fatality_Ratio" {
		for ind, line := range lines {
			if ind > 0 {
				line = strings.Replace(line, "'", "''", -1)
				ws := splitLine(line, []int{0, 8, 9, 10})
				record := db.Record{
					Date:               Date,
					FIPS:               ws[0],
					Admin2:             ws[1],
					Province:           ws[2],
					Country:            ws[3],
					LastUpdate:         ws[4],
					Lat:                ws[5],
					Long:               ws[6],
					Confirmed:          ws[7],
					Deaths:             ws[8],
					Recovered:          ws[9],
					Active:             ws[10],
					CombinedKey:        ws[11],
					IncidentRate:       ws[12],
					CaseFatalityRation: ws[13],
				}
				if len(record.Province) == 0 {
					record.Province = "__total__"
				} else {
					if _, ok := records[record.Country+"Confirmed"]; ok {
						records[record.Country+"Confirmed"] += Datas.AtoI(record.Confirmed)
						records[record.Country+"Deaths"] += Datas.AtoI(record.Deaths)
						records[record.Country+"Recovered"] += Datas.AtoI(record.Recovered)
						records[record.Country+"Active"] += Datas.AtoI(record.Active)
					} else {
						recordsCountry = append(recordsCountry, record.Country)
						records[record.Country+"Confirmed"] = Datas.AtoI(record.Confirmed)
						records[record.Country+"Deaths"] = Datas.AtoI(record.Deaths)
						records[record.Country+"Recovered"] = Datas.AtoI(record.Recovered)
						records[record.Country+"Active"] = Datas.AtoI(record.Active)
					}
				}
				database.Insert(record)
			}
		}
	} else if lines[0] == "FIPS,Admin2,Province_State,Country_Region,Last_Update,Lat,Long_,Confirmed,Deaths,Recovered,Active,Combined_Key" {
		for ind, line := range lines {
			if ind > 0 {
				line = strings.Replace(line, "'", "''", -1)
				ws := splitLine(line, []int{0, 8, 9, 10})
				record := db.Record{
					Date:        Date,
					FIPS:        ws[0],
					Admin2:      ws[1],
					Province:    ws[2],
					Country:     ws[3],
					LastUpdate:  ws[4],
					Lat:         ws[5],
					Long:        ws[6],
					Confirmed:   ws[7],
					Deaths:      ws[8],
					Recovered:   ws[9],
					Active:      ws[10],
					CombinedKey: ws[11],
				}
				if len(record.Province) == 0 {
					record.Province = "__total__"
				} else {
					if _, ok := records[record.Country+"Confirmed"]; ok {
						records[record.Country+"Confirmed"] += Datas.AtoI(record.Confirmed)
						records[record.Country+"Deaths"] += Datas.AtoI(record.Deaths)
						records[record.Country+"Recovered"] += Datas.AtoI(record.Recovered)
						records[record.Country+"Active"] += Datas.AtoI(record.Active)
					} else {
						recordsCountry = append(recordsCountry, record.Country)
						records[record.Country+"Confirmed"] = Datas.AtoI(record.Confirmed)
						records[record.Country+"Deaths"] = Datas.AtoI(record.Deaths)
						records[record.Country+"Recovered"] = Datas.AtoI(record.Recovered)
						records[record.Country+"Active"] = Datas.AtoI(record.Active)
					}
				}
				database.Insert(record)
			}
		}
	} else if strings.Index(lines[0], "Province/State,Country/Region,Last Update,Confirmed,Deaths,Recovered") != -1 {
		for ind, line := range lines {
			if ind > 0 {
				line = strings.Replace(line, "'", "''", -1)
				ws := splitLine(line, []int{3, 4, 5})
				ws[5] = strings.Split(ws[5], ".")[0]
				record := db.Record{
					Date:       Date,
					Province:   ws[0],
					Country:    ws[1],
					LastUpdate: ws[2],
					Confirmed:  ws[3],
					Deaths:     ws[4],
					Recovered:  ws[5],
				}
				if len(record.Province) == 0 {
					record.Province = "__total__"
				} else {
					if _, ok := records[record.Country+"Confirmed"]; ok {
						records[record.Country+"Confirmed"] += Datas.AtoI(record.Confirmed)
						records[record.Country+"Deaths"] += Datas.AtoI(record.Deaths)
						records[record.Country+"Recovered"] += Datas.AtoI(record.Recovered)
					} else {
						recordsCountry = append(recordsCountry, record.Country)
						records[record.Country+"Confirmed"] = Datas.AtoI(record.Confirmed)
						records[record.Country+"Deaths"] = Datas.AtoI(record.Deaths)
						records[record.Country+"Recovered"] = Datas.AtoI(record.Recovered)
						records[record.Country+"Active"] = 0
					}
				}
				database.Insert(record)
			}
		}
	} else if strings.Index(lines[0], "Province/State,Country/Region,Last Update,Confirmed,Deaths,Recovered,Latitude,Longitude") != -1 {
		for ind, line := range lines {
			if ind > 0 {
				line = strings.Replace(line, "'", "''", -1)
				ws := splitLine(line, []int{3, 4, 5})
				record := db.Record{
					Date:       Date,
					Province:   ws[0],
					Country:    ws[1],
					LastUpdate: ws[2],
					Confirmed:  ws[3],
					Deaths:     ws[4],
					Recovered:  ws[5],
					Lat:        ws[6],
					Long:       ws[7],
				}
				if len(record.Province) == 0 {
					record.Province = "__total__"
				} else {
					if _, ok := records[record.Country+"Confirmed"]; ok {
						records[record.Country+"Confirmed"] += Datas.AtoI(record.Confirmed)
						records[record.Country+"Deaths"] += Datas.AtoI(record.Deaths)
						records[record.Country+"Recovered"] += Datas.AtoI(record.Recovered)
					} else {
						recordsCountry = append(recordsCountry, record.Country)
						records[record.Country+"Confirmed"] = Datas.AtoI(record.Confirmed)
						records[record.Country+"Deaths"] = Datas.AtoI(record.Deaths)
						records[record.Country+"Recovered"] = Datas.AtoI(record.Recovered)
						records[record.Country+"Active"] = 0
					}
				}
				database.Insert(record)
			}
		}
	} else {
		fmt.Println("排除文件 " + filepath.Base(csvPath))
	}
	for _, country := range recordsCountry {
		record := db.Record{
			Date:      Date,
			Province:  "__total__",
			Country:   country,
			Confirmed: strconv.Itoa(records[country+"Confirmed"]),
			Deaths:    strconv.Itoa(records[country+"Deaths"]),
			Recovered: strconv.Itoa(records[country+"Recovered"]),
			Active:    strconv.Itoa(records[country+"Active"]),
		}
		database.Insert(record)
	}
}
