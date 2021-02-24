package SqliteSql

import (
	"database/sql"
)

/*
使用方法
sql := "select note,count(*) from 表名 GROUP BY 分组字段"
entities := QueryGroupCountStringInt(sql,db)
for _,entity := range entities { ... }
*/
type GroupCountStringInt struct {
	Text  string
	Count int
}

func QueryGroupCountStringInt(sqlStr string, db *sql.DB) []GroupCountStringInt {
	ParseRowsGroupCountStringInt := func(rows *sql.Rows) []GroupCountStringInt {
		groups := []GroupCountStringInt{}
		for rows.Next() {
			group := GroupCountStringInt{}
			rows.Scan(
				&group.Text,
				&group.Count)
			groups = append(groups, group)
		}
		return groups
	}
	if rows, err := db.Query(sqlStr); err == nil {
		return ParseRowsGroupCountStringInt(rows)
	} else {
		// todo 错误暂时无解决方案
		panic(err)
	}
}

func QueryCount(table string, db *sql.DB) int {
	getCount := func(rows *sql.Rows) int {
		var count int
		for rows.Next() {
			rows.Scan(&count)
			break
		}
		return count
	}
	if rows, err := db.Query(`select count(*) from ` + table); err == nil {
		return getCount(rows)
	} else {
		// todo 错误暂时无解决方案
		panic(err)
	}
}

// 执行无返回的数据库查询
func ExecSqlString(db *sql.DB, sqlstring string) {
	//fmt.Print(sqlstring)
	statement, err := db.Prepare(sqlstring)
	if err != nil {
		panic(err)
	}
	defer statement.Close()
	statement.Exec()
}
