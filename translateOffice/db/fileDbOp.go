package db

import (
	"database/sql"
	"path"
	"public.sunibas.cn/go_public/utils/SqliteSql"
	"strings"
)

type FileDb struct {
	Path   string
	FileDb *sql.DB
}

var FileTableName = "file"

type File struct {
	Id      string
	Content string
}

func (db *FileDb) Init() {
	db.Path = "_tmp_file_.db"
	if len(dbPath) != 0 {
		db.Path = path.Join(dbPath, "_tmp_file_.db")
	}
	var err error
	db.FileDb, err = sql.Open("sqlite3", db.Path)
	if err != nil {
		panic(err)
	}
	statement, _ := db.FileDb.Prepare(`CREATE TABLE IF NOT EXISTS ` + FileTableName + ` (
		"id" TEXT,
		"content" TEXT,
		PRIMARY KEY ("id")
	);`)
	statement.Exec()
	statement, _ = db.FileDb.Prepare(`CREATE INDEX IF NOT EXISTS 'dir_id' ON '` + FileTableName + `' ( 'id' ASC );`)
	statement.Exec()
}

func (db *FileDb) InsertFile(id, content string) {
	SqliteSql.ExecSqlString(db.FileDb, "insert into "+FileTableName+"(`id`,`content`) values('"+id+"','"+content+"')")
}

func (db *FileDb) InsertValues(id, content []string) {
	SqliteSql.ExecSqlString(db.FileDb, "insert into "+FileTableName+"(`id`,`name`,`content`) values"+
		toFileValues(id, content))
}

func toFileValues(id, content []string) string {
	values := []string{}
	for i := 0; i < len(id); i++ {
		//values = append(
		//	values,"('" + id[i] + "','" + fileid + "','" + type_ + "','" + strings.Replace(content[i],"'","''",-1) + "')")
		values = append(values, "('"+strings.Join([]string{
			id[i], content[i],
		}, "','")+"')")
	}
	return strings.Join(values, ",")
}

func (db *FileDb) SelectByIdFromFile(id string) []File {
	sql := `select id,content from ` + FileTableName + ` where id='` + id + `' `

	if rows, err := db.FileDb.Query(sql); err != nil {
		panic(err)
	} else {
		dirs := []File{}
		for rows.Next() {
			dir := File{}
			rows.Scan(
				&dir.Id,
				&dir.Content,
			)
			dirs = append(dirs, dir)
		}
		return dirs
	}
}
