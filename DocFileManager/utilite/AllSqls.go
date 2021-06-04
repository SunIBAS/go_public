package utilite

import (
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"strings"
)

// 跟目录 id
const RootId = "0"

// 原文符号 org
const OrgType = "org"

const (
	FileStatus_add      = "add"
	FileStatus_check    = "check"
	FileStatus_read     = "read"
	FileStatus_readover = "readover"
)

type DirDir struct {
	Id         string
	Parent     string
	Filename   string
	Type       string
	Createdate string
	Lang       string
	Fileext    string
	Absolute   string
	Note       string
	Status     string
}
type DirTheStatus struct {
	Id   string
	Sta  string
	Bak1 string
	Bak2 string
}

func CreateDirTable(dirDb *sql.DB) {
	// dir 用于记录文件夹的层级关系和文件的基本信息， lang 字段可以用于记录文件名的语言条目
	//		absolute 文件的绝对路径
	//		status 文件处理状态，分别为 add(仅加入到文件结构中) check(完成文件检查) read（开始处理文件）readover(处理完成)
	// lang 中的 type = org 时表示记录的内容为原始名称
	// theStatus 用于记录软件运行状态信息
	sqls := []string{
		`CREATE TABLE "dir" (
			"id" TEXT,
			"parent" TEXT,
			"filename" TEXT,
			"type" TEXT,
			"createdate" INTEGER,
			"lang" TEXT,
			"fileext" TEXT,
			"absolute" TEXT,
			"note" TEXT,
			"status" TEXT,
			PRIMARY KEY ("id")
		);`,
		`CREATE INDEX 'dir_id' ON 'dir' ( 'id' ASC );`,
		`CREATE TABLE "lang" (
			"id" TEXT,
			"type" TEXT,
			"content" TEXT,
			PRIMARY KEY ("id", "type")
		);`,
		`CREATE INDEX 'lang_id' ON 'lang' ( 'id' ASC );`,
		`insert into dir values("0",null,"root","dir","","0",null)`,
		`CREATE TABLE "theStatus" (
			"id" TEXT,
			"sta" TEXT,
			"bak1" TEXT,
			"bak2" TEXT,
			PRIMARY KEY ("id")
		);`,
	}

	for _, sql := range sqls {
		statement, _ := dirDb.Prepare(sql)
		if statement != nil {
			statement.Exec()
		}
	}
}
func CreateLangTable(langDb *sql.DB) {
	// fileTemplate 单纯关联 dir 数据库中的 dir_id 和 文件模板
	// lang 用于记录翻译或原文 type=org 表示原文
	sqls := []string{
		`CREATE TABLE "fileTemplate" (
			"id" TEXT,
			"fileid" TEXT,
			"filecontent" BLOB,
			"bak1" TEXT,
			"bak2" TEXT,
			PRIMARY KEY ("id")
		);`,
		`CREATE INDEX "ft_fid"
		ON "fileTemplate" (
			"fileid" ASC
		);`,
		`CREATE TABLE "lang" (
			"id" TEXT,
			"fileid" TEXT,
			"type" TEXT,
			"content" TEXT,
			PRIMARY KEY ("id", "type")
		);`,
		`CREATE INDEX 'lang_id' ON 'lang' ( 'id' ASC );`,
		`CREATE INDEX 'lang_fid' ON 'lang' ( 'fileid' ASC );`,
	}

	for _, sql := range sqls {
		statement, _ := langDb.Prepare(sql)
		if statement != nil {
			statement.Exec()
		}
	}
}
func CreateId() string {
	var err error
	return uuid.Must(uuid.NewV4(), err).String()
}
func Sql_insert_into_lang(id, type_, content string) string {
	var s = []string{
		`insert into lang(id,type,content) values('`,
		id,
		`','`,
		type_,
		`','`,
		content,
		`');`,
	}
	return strings.Join(s, "")
}
func Sql_insert_into_dir(id, parent, filename, type_, createdate, lang, fileext, absolute, note, status string) string {
	var s = []string{
		`insert into lang(id,parent,filename,type,createdate,lang,fileext,absolute,note,status) values('`,
		id,
		`','`,
		parent,
		`','`,
		filename,
		`','`,
		type_,
		`','`,
		createdate,
		`','`,
		lang,
		`','`,
		fileext,
		`','`,
		absolute,
		`','`,
		note,
		`','`,
		status,
		`');`,
	}
	return strings.Join(s, "")
}
func Sql_insert_into_status(id, sta, bak1, bak2 string) string {
	var s = []string{
		`insert into theStatus(id,sta,bak1,bak2) values('`,
		id,
		`','`,
		sta,
		`','`,
		bak1,
		`','`,
		bak2,
		`');`,
	}
	return strings.Join(s, "")
}

// where = "" 或 where = " where id='1'"
func Sql_select_from_dir(dirDb *sql.DB, where string) []DirDir {
	sql := `select id,parent,filename,type createdate,lang,fileext,absolute,note,status from dir ` + where
	if rows, err := dirDb.Query(sql); err == nil {
		dirs := []DirDir{}
		for rows.Next() {
			dir := DirDir{}
			rows.Scan(
				&dir.Id,
				&dir.Parent,
				&dir.Filename,
				&dir.Type,
				&dir.Createdate,
				&dir.Lang,
				&dir.Fileext,
				&dir.Absolute,
				&dir.Note,
				&dir.Status,
			)
			dirs = append(dirs, dir)
		}
		return dirs
	} else {
		// todo 错误暂时无解决方案
		panic(err)
	}
}
func Sql_select_from_the_status(dirDb *sql.DB, where string) []DirTheStatus {
	sql := `select id,sta,bak1,bak2 from theStatus ` + where
	if rows, err := dirDb.Query(sql); err == nil {
		status := []DirTheStatus{}
		for rows.Next() {
			statu := DirTheStatus{}
			rows.Scan(
				&statu.Id,
				&statu.Sta,
				&statu.Bak1,
				&statu.Bak2,
			)
			status = append(status, statu)
		}
		return status
	} else {
		// todo 错误暂时无解决方案
		panic(err)
	}
}
