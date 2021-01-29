package Times

import (
	"public.sunibas.cn/go_public/utils/Datas"
	"time"
)

var day = 3600 * 24

// fr = 20200101 to = 20200102 => 1
func Between(fr, to string) int64 {
	frMonth := time.Month(Datas.AtoI(fr[4:6]))
	frTime := time.Date(Datas.AtoI(fr[:4]), frMonth, Datas.AtoI(fr[6:8]), 0, 0, 0, 0, time.Local).Unix()
	toMonth := time.Month(Datas.AtoI(to[4:6]))
	toTime := time.Date(Datas.AtoI(to[:4]), toMonth, Datas.AtoI(to[6:8]), 0, 0, 0, 0, time.Local).Unix()
	return (toTime - frTime) / int64(day)
}

// fr = 20200101 days = 1 => 20200102
// format = 1 20200102
// format = 2 2020-01-02
func DaysToDay(fr string, days, format int) string {
	frMonth := time.Month(Datas.AtoI(fr[4:6]))
	frTime := time.Date(Datas.AtoI(fr[:4]), frMonth, Datas.AtoI(fr[6:8]), 0, 0, 0, 0, time.Local).Unix()
	frTime += int64(days * day)
	if format == 1 {
		return time.Unix(frTime, 0).Format("20060102")
	} else {
		return time.Unix(frTime, 0).Format("2006-01-02")
	}
}

// Formater = "2006-01-02 15:04:05" 日期格式
type DaysNexter struct {
	Date     time.Time
	ToDate   string
	Time     string
	Formater string
}

func (dn *DaysNexter) Init(year, month, date int) {
	dn.Date = time.Date(year, time.Month(month), date, 0, 0, 0, 0, time.UTC)
	dn.ToDate = time.Now().Format(dn.Formater)
}
func (dn *DaysNexter) NextDate() string {
	dn.Date = dn.Date.AddDate(0, 0, 1)
	return dn.Date.Format(dn.Formater)
}

// 循环到今天
func (dn *DaysNexter) NextDateToTD() (string, bool) {
	d := dn.NextDate()
	return d, d == dn.ToDate
}
