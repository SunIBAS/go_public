package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"public.sunibas.cn/go_public/utils/Network"
	"public.sunibas.cn/go_public/utils/SqliteSql"

	//"sunibas.cn/go_utils/utils/DirAndFile"
	//"sunibas.cn/go_utils/utils/Network"
	//"sunibas.cn/go_utils/utils/SqliteSql"
	"sync"
	"time"
)

var wg sync.WaitGroup //定义一个同步等待的组
var _database *sql.DB

func help() {
	fmt.Println("dl_btc.exe req null.db")
	fmt.Println("dl_btc.exe json null.db")
}

func main() {
	if len(os.Args) > 2 {
		var err error
		_database, err = sql.Open("sqlite3", os.Args[2])
		if err != nil {
			panic(err)
		}
		defer _database.Close()
		if os.Args[1] == "req" {
			if len(os.Args) == 3 {
				requestDataAndSaveToSqlite()
			}
		} else if os.Args[1] == "json" {
			toJson()
		} else {
			help()
		}
	} else {
		help()
	}
}

type record struct {
	Id      string
	Content string
	Time    string
}

func toJson() {
	if rows, err := _database.Query(`select id,content,time from record ORDER BY id`); err == nil {
		records := []record{}
		for rows.Next() {
			rec := record{}
			rows.Scan(
				&rec.Id,
				&rec.Content,
				&rec.Time,
			)
			records = append(records, rec)
		}
		if content, err := json.Marshal(records); err == nil {
			DirAndFile.WriteWithIOUtilByte("out.json", content)
		} else {
			panic(err)
		}

	} else {
		// todo 错误暂时无解决方案
		panic(err)
	}
}

func requestDataAndSaveToSqlite() {
	//_database,err := sql.Open("sqlite3", "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\_test\\1to100\\null.db")
	for {
		time.Sleep(time.Duration(time.Second * 2))
		str, _ := Network.Get("https://www.coinhako.com/api/v1/price/currency/BTCUSD")
		SqliteSql.ExecSqlString(_database, "insert into record('content','time') values('"+str+"','"+time.Now().Format("2006-01-02 15:04:05")+"')")
	}
}
