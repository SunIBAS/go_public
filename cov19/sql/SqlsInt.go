package sqlsint

import (
	//"database/sql"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type SqlsInt struct {
	DbPath string
	Database * sql.DB
	ShowSql bool
}

type Table struct {
	Id string `id`
	CityFrom string `city_from`
	CityTo string `city_to`
	Num string `num`
	Date string `date`
	Whos string `whos`
	Note string `note`
}
type NameTable struct {
	Name string `json:"name"`
}

func (si * SqlsInt)Init() {
	var err error
	si.Database,err = sql.Open("sqlite3", si.DbPath)
	if err != nil {
		panic(err)
	}
}

func (si SqlsInt)Close() {
	si.Database.Close()
}

/**
 * from 	来源	* 表示发生在当地而非输入
 * to 		省市
 * num 		数量
 * date		日期
 * note		说明内容
 * itype	写入到哪个表，
 *		daily_inj	每日输入
 *		total_inj	总输入
 *		daily_wzz	每日无症状
 *		total_wzz	总无症状
 */
func (si SqlsInt)Insert(from,to,date,num,itype,whos,note string)  {
	values := []string{
		from,to,num,date,whos,note,
	}
	insertSql := `insert into ` + itype + `(city_from,city_to,num,date,whos,note) values('` + strings.Join(values,"','") +  `')`
	if si.ShowSql {
		fmt.Println(insertSql)
	}
	statement,_ := si.Database.Prepare(insertSql)
	statement.Exec()
}

func (si SqlsInt)SelectD(date,itype string)[]Table {
	where := []string{
		`date='`,
		date,
		"'",
	}
	return si.Select(strings.Join(where,""),itype)
}

func (si SqlsInt)SelectDI(from,date,itype string) []Table {
	where := []string{
		`city_from='`,
		from,
		`' and date='`,
		date,
		"'",
	}
	return si.Select(strings.Join(where,""),itype)
}

func (si SqlsInt)Select(where,itype string) []Table {
	selectSql := []string {
		`select id,city_from,city_to,num,date,whos,note from `,
		itype,
		` where `,
		where,
	}
	sql := strings.Join(selectSql,"")
	if si.ShowSql {
		fmt.Println(sql)
	}
	if rows,err := si.Database.Query(sql);
		err == nil {
			return rowToEntity(rows)
	} else {
			//panic(err)
		return []Table{}
	}
}

func rowToEntity(rows * sql.Rows) (tables []Table) {
	for rows.Next() {
		table := Table{}
		rows.Scan(
			&table.Id,
			&table.CityFrom,
			&table.CityTo,
			&table.Num,
			&table.Date,
			&table.Whos,
			&table.Note)
		tables = append(tables, table)
	}
	return
}

func (si SqlsInt)SelectNames()(tables []NameTable) {
	sql := `select name from names`
	if si.ShowSql {
		fmt.Println(sql)
	}
	if rows,err := si.Database.Query(sql);
		err == nil {
		for rows.Next() {
			table := NameTable{}
			rows.Scan(&table.Name)
			tables = append(tables, table)
		}
	}
	return tables
}
