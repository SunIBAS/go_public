package db

import (
	"database/sql"
	"path"
	"public.sunibas.cn/go_public/utils/SqliteSql"
	"strings"
)

type LangDb struct {
	Path   string
	LangDb *sql.DB
}

var (
	WorkBook  string = "workboox"
	SheetName string = "sheetname"
	Word      string = "word"
	XLSX      string = "xlsx"
	DOCX      string = "docx"
)

type Lang struct {
	Id      string
	Fileid  string
	Type    string
	Content string
	Tran    string
}

func (db *LangDb) Init() {
	db.Path = "_tmp_lang_.db"
	if len(dbPath) != 0 {
		db.Path = path.Join(dbPath, "_tmp_lang_.db")
	}
	var err error
	db.LangDb, err = sql.Open("sqlite3", db.Path)
	if err != nil {
		panic(err)
	}
	statement, _ := db.LangDb.Prepare(`CREATE TABLE IF NOT EXISTS lang(
		"id" TEXT,
		"fileid" TEXT,
		"type" TEXT,
		"content" TEXT,
		"tran" TEXT,
		PRIMARY KEY ("id", "type")
	);`)
	statement.Exec()
	statement, _ = db.LangDb.Prepare(`CREATE INDEX IF NOT EXISTS 'lang_id' ON 'lang' ( 'id' ASC );`)
	statement, _ = db.LangDb.Prepare(`CREATE INDEX IF NOT EXISTS 'lang_fileid' ON 'lang' ( 'fileid' ASC );`)
	statement.Exec()
}

func (db *LangDb) Insert(id, fileid, type_, content string) {
	SqliteSql.ExecSqlString(db.LangDb, "insert into lang(`id`,`fileid`,`type`,`content`) values('"+
		id+"','"+fileid+"','"+type_+"','"+strings.Replace(content, "'", "''", -1)+"')")
}

func (db *LangDb) InsertValues(id []string, fileid, type_ string, content []string) {
	SqliteSql.ExecSqlString(db.LangDb, "insert into lang(`id`,`fileid`,`type`,`content`) values"+langToValues(id, fileid, type_, content))
}

func langToValues(id []string, fileid, type_ string, content []string) string {
	values := []string{}
	for i := 0; i < len(id); i++ {
		values = append(values, "('"+id[i]+"','"+fileid+"','"+type_+"','"+strings.Replace(content[i], "'", "''", -1)+"')")
	}
	return strings.Join(values, ",")
}

func (db *LangDb) Select(fileid string, limit string) []Lang {
	sql := `select * from lang where fileid='` + fileid + `' `
	if fileid == "" {
		sql = `select * from lang `
	}
	if len(limit) > 0 {
		sql += limit
	}
	if rows, err := db.LangDb.Query(sql); err != nil {
		panic(err)
	} else {
		langs := []Lang{}
		for rows.Next() {
			lang := Lang{}
			rows.Scan(
				&lang.Id,
				&lang.Fileid,
				&lang.Type,
				&lang.Content,
				&lang.Tran,
			)
			langs = append(langs, lang)
		}
		return langs
	}
}

//func (db * LangDb)Clear() {
//	db.LangDb.Exec("delete from lang where 1=1")
//}

func (db *LangDb) UpdateTran(id, tran string) {
	db.LangDb.Exec(`update lang set tran='` + tran + `' where id='` + id + `'`)
}

func (db *LangDb) UpdateTrans(id, tran []string) {
	sqls := []string{}
	for ind := 0; ind < len(id); ind++ {
		sqls = append(sqls, `update lang set tran='`+tran[ind]+`' where id='`+id[ind]+`';`)
	}
	db.LangDb.Exec(strings.Join(sqls, ";"))
}
