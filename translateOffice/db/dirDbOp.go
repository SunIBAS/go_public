package db

import (
	"database/sql"
	"path"
	"public.sunibas.cn/go_public/utils/SqliteSql"
	"strings"
)

type DirDb struct {
	Path  string
	DirDb *sql.DB
}

var (
	DirTableName = "dir"
	DirString    = "dir"
	FileString   = "file"
)

type Dir struct {
	Id       string
	Name     string
	RepName  string
	Type     string
	ParentId string
	Tag      string
}

func (db *DirDb) Init() {
	db.Path = "_tmp_dir_.db"
	if len(dbPath) != 0 {
		db.Path = path.Join(dbPath, "_tmp_dir_.db")
	}
	var err error
	db.DirDb, err = sql.Open("sqlite3", db.Path)
	if err != nil {
		panic(err)
	}
	statement, _ := db.DirDb.Prepare(`CREATE TABLE IF NOT EXISTS ` + DirTableName + ` (
		"id" TEXT,
		"name" TEXT,
		"repname" TEXT,
		"type" TEXT,
		"parentid" TEXT,
		"tag" TEXT,
		PRIMARY KEY ("id", "type")
	);`)
	statement.Exec()
	statement, _ = db.DirDb.Prepare(`CREATE INDEX IF NOT EXISTS 'dir_id' ON '` + DirTableName + `' ( 'id' ASC );`)
	statement.Exec()
}

func (db *DirDb) InsertDir(id, name, repname, type_, parentid, tag string) {
	SqliteSql.ExecSqlString(db.DirDb, "insert into "+DirTableName+"(`id`,`name`,`repname`,`type`,`parentid`,`tag`) values('"+
		id+"','"+name+"','"+repname+"','"+type_+"','"+parentid+"','"+tag+"')")
}

func (db *DirDb) InsertValues(id, name, repname, type_, parentid []string, tag string) {
	SqliteSql.ExecSqlString(db.DirDb, "insert into "+DirTableName+"(`id`,`name`,`repname`,`type`,`parentid`,`tag`) values"+
		toDirValues(id, name, repname, type_, parentid, tag))
}

func toDirValues(id, name, repname, type_, parentid []string, tag string) string {
	values := []string{}
	for i := 0; i < len(id); i++ {
		//values = append(
		//	values,"('" + id[i] + "','" + fileid + "','" + type_ + "','" + strings.Replace(content[i],"'","''",-1) + "')")
		values = append(values, "('"+strings.Join([]string{
			id[i], name[i], repname[i], type_[i], parentid[i], tag,
		}, "','")+"')")
	}
	return strings.Join(values, ",")
}

func (db *DirDb) SelectByTagFromDir(tag string) []Dir {
	sql := `select id,name,repname,type,parentid,tag from ` + DirTableName + ` where tag='` + tag + `' `

	if rows, err := db.DirDb.Query(sql); err != nil {
		panic(err)
	} else {
		dirs := []Dir{}
		for rows.Next() {
			dir := Dir{}
			rows.Scan(
				&dir.Id,
				&dir.Name,
				&dir.RepName,
				&dir.Type,
				&dir.ParentId,
				&dir.Tag,
			)
			dirs = append(dirs, dir)
		}
		return dirs
	}
}

func (db *DirDb) SelectByParentIdFromDir(parentid string) []Dir {
	sql := `select id,name,repname,type,parentid,tag from ` + DirTableName + ` where parentid='` + parentid + `' `

	if rows, err := db.DirDb.Query(sql); err != nil {
		panic(err)
	} else {
		dirs := []Dir{}
		for rows.Next() {
			dir := Dir{}
			rows.Scan(
				&dir.Id,
				&dir.Name,
				&dir.RepName,
				&dir.Type,
				&dir.ParentId,
				&dir.Tag,
			)
			dirs = append(dirs, dir)
		}
		return dirs
	}
}

//func (db * DirDb)Clear() {
//	db.DirDb.Exec("delete from " + DirTableName + " where 1=1")
//}

//func (db * DirDb)UpdateTran(id,tran string) {
//	db.DirDb.Exec(`update ` + DirTableName + ` set tran='` + tran + `' where id='` + id + `'`)
//}
