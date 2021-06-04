package db

import (
	"database/sql"
	"public.sunibas.cn/go_public/utils/SqliteSql"
	"strings"
)

type Db struct {
	Path string
	Db   *sql.DB
}

type Record struct {
	Id                 string
	Date               string
	FIPS               string
	Admin2             string
	Province           string
	Country            string
	LastUpdate         string
	Lat                string
	Long               string
	Confirmed          string
	Deaths             string
	Recovered          string
	Active             string
	CombinedKey        string
	IncidentRate       string
	CaseFatalityRation string
}

func (db *Db) Init() {
	db.Path = "tmp.db"
	var err error
	db.Db, err = sql.Open("sqlite3", db.Path)
	if err != nil {
		panic(err)
	}
	statement, _ := db.Db.Prepare(`CREATE TABLE IF NOT EXISTS record(
		"id" TEXT,
		"date" BIGINT,
		"fips" INT,
		"admin2" TEXT,
		"province" TEXT,
		"country" TEXT,
		"lastUpdate" TEXT,
		"lat" FLOAT,
		"long" FLOAT,
		"confirmed" BIGINT,
		"deaths" BIGINT,
		"recovered" BIGINT,
		"active" BIGINT,
		"combinekey" TEXT,
		"incidentrate" FLOAT,
		"casefatalityration" FLOAT,
		PRIMARY KEY ("id")
	);`)
	statement.Exec()
	statement, _ = db.Db.Prepare(`CREATE INDEX IF NOT EXISTS 'record_id' ON 'record' ( 'id' ASC );`)
	//statement,_ = db.Db.Prepare(`CREATE INDEX IF NOT EXISTS 'record_fileid' ON 'record' ( 'fileid' ASC );`)
	statement.Exec()
}

func (db *Db) Insert(record Record) {
	if record.Country == "Czechia" {
		record.Country = "Czech Republic"
	}
	if record.Country == "Korea, South" {
		record.Country = "South Korea"
	}
	record.Id = record.Date + record.Country + record.Province + record.Admin2 + record.FIPS
	fields := []string{
		`id`, `date`, `fips`, `admin2`, `province`, `country`, `lastUpdate`, `lat`, `long`,
		`confirmed`, `deaths`, `recovered`, `active`, `combinekey`,
		`incidentrate`, `casefatalityration`,
	}
	values := []string{
		record.Id, record.Date, record.FIPS, record.Admin2, record.Province, record.Country, record.LastUpdate,
		record.Lat, record.Long, record.Confirmed, record.Deaths, record.Recovered, record.Active,
		record.CombinedKey, record.IncidentRate, record.CaseFatalityRation,
	}
	//updateStr := strings.Join([]string{
	//	"confirmed='" + record.Confirmed + "'",
	//	"deaths='" + record.Deaths + "'",
	//	"recovered='" + record.Recovered + "'",
	//},",")

	sql := `INSERT OR REPLACE INTO record (` + strings.Join(fields, ",") + `)
values('` + strings.Join(values, "','") + `')`
	//fmt.Println(sql)
	SqliteSql.ExecSqlString(db.Db, sql)
}

func (db *Db) Select(country string, province string) []Record {
	sql := `select * from record where country="` + country + `" and province="` + province + `" order by date asc`
	if rows, err := db.Db.Query(sql); err != nil {
		panic(err)
	} else {
		records := []Record{}
		for rows.Next() {
			record := Record{}
			rows.Scan(
				&record.Id,
				&record.Date,
				&record.FIPS,
				&record.Admin2,
				&record.Province,
				&record.Country,
				&record.LastUpdate,
				&record.Lat,
				&record.Long,
				&record.Confirmed,
				&record.Deaths,
				&record.Recovered,
				&record.Active,
				&record.CombinedKey,
				&record.IncidentRate,
				&record.CaseFatalityRation,
			)
			records = append(records, record)
		}
		return records
	}
}

func (db *Db) Clear() {
	db.Db.Exec("delete from record where 1=1")
}

func (db *Db) UpdateTran(id, tran string) {
	db.Db.Exec(`update record set tran='` + tran + `' where id='` + id + `'`)
}
