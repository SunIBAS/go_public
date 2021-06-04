package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os"
	"public.sunibas.cn/go_public/utils/Datas"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"public.sunibas.cn/go_public/utils/Network"
	"public.sunibas.cn/go_public/utils/SqliteSql"
	"strconv"

	//"sunibas.cn/go_utils/utils/DirAndFile"
	//"sunibas.cn/go_utils/utils/Network"
	//"sunibas.cn/go_utils/utils/SqliteSql"
	"sync"
	"time"
)

var wg sync.WaitGroup //定义一个同步等待的组
var _database *sql.DB
var port string
var lastData string

func help() {
	fmt.Println("请求数据并记录到数据库中")
	fmt.Println("dl_btc.exe req null.db")
	fmt.Println("")
	fmt.Println("将数据库中所有内容写入到 json 文件中")
	fmt.Println("dl_btc.exe json null.db")
	fmt.Println("")
	fmt.Println("仅开启查询服务（服务只能用 post，可跨域）")
	fmt.Println("dl_btc.exe server null.db 8080")
	fmt.Println(`curl -d "{\"Content\": \"100\", \"Method\": \"limit\"}" -H 'Content-Type: application/json' http://localhost:8087/api`)
	fmt.Println(`fetch('http://localhost:8080/api',{
    headers: {
        'Content-Type': 'application/json;charset=UTF-8'
    },
    method: 'post',
    body: JSON.stringify({Content:'100',Method:'limit'})
}).then(_ => _.text()).then(JSON.parse).then(_ => JSON.parse(_.Content)).then(d => d.map(d => { d.Content = JSON.parse(d.Content);return d; })).then(ret => window.ret = ret);`)
	fmt.Println("")
	fmt.Println("开启记录同时开启服务")
	fmt.Println("dl_btc.exe qs null.db 8080")
	fmt.Println("")
	fmt.Println("删除相同的数据")
	fmt.Println("dl_btc.exe removes null.db")
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
			requestDataAndSaveToSqlite()
		} else if os.Args[1] == "json" {
			toJson()
		} else if os.Args[1] == "server" {
			port = os.Args[3]
			server()
		} else if os.Args[1] == "qs" {
			port = os.Args[3]
			go requestDataAndSaveToSqlite()
			server()
		} else if os.Args[1] == "removes" {
			removeSame()
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

// `select id,content,time from record ORDER BY id`
// sqlTail = 'ORDER BY id'
func queryRecord(sqlTail string) []record {
	records := []record{}
	if rows, err := _database.Query(`select id,content,time from record ` + sqlTail); err == nil {
		for rows.Next() {
			rec := record{}
			rows.Scan(
				&rec.Id,
				&rec.Content,
				&rec.Time,
			)
			records = append(records, rec)
		}
	}
	return records
}
func toJson() {
	if content, err := json.Marshal(queryRecord(`ORDER BY id`)); err == nil {
		DirAndFile.WriteWithIOUtilByte("out.json", content)
	} else {
		panic(err)
	}
}

func requestDataAndSaveToSqlite() {
	//_database,err := sql.Open("sqlite3", "C:\\Users\\IBAS\\go\\src\\sunibas.cn\\go_utils\\_test\\1to100\\null.db")
	for {
		time.Sleep(time.Duration(time.Second * 10))
		str, _ := Network.Get("https://www.coinhako.com/api/v1/price/currency/BTCUSD")
		if lastData != str {
			lastData = str
			SqliteSql.ExecSqlString(_database, "insert into record('content','time') values('"+str+"','"+time.Now().Format("2006-01-02 15:04:05")+"')")
		}
	}
}

func server() {
	http.HandleFunc("/api", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		writer.Header().Set("content-type", "application/json")             //返回数据格式是json
		var rObj = Network.ReturnObj{
			Code:    0,
			Message: "",
			Content: "",
		}
		if p, err := Network.PostParams(request); err == nil {
			if p.Method == "limit" {
				rObj.SetContent(queryRecord(` ORDER BY id desc limit ` + p.Content))
			} else if p.Method == "insert" {
				SqliteSql.ExecSqlString(_database, `insert into record(content,time) values(`+p.Content+`);`)
				rObj.SetSuccess("ok")
			}
		} else {
			rObj.SetFail(err.Error())
		}
		rObj.Send(writer)
	})
	srv := &http.Server{
		Addr:           `:` + port,
		Handler:        nil,
		ReadTimeout:    time.Duration(5) * time.Minute,
		WriteTimeout:   time.Duration(5) * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("服务已开启 http://localhost:" + port)
}

func removeSame() {
	var maxId, minId int
	var _last string
	rec := queryRecord(`order BY id desc limit 1`)
	maxId = Datas.AtoI(rec[0].Id)
	rec = queryRecord(`order by id asc limit 1`)
	minId = Datas.AtoI(rec[0].Id)
	fmt.Println(`minId = ` + strconv.Itoa(minId) + ` maxId = ` + strconv.Itoa(maxId))
	for ; minId < maxId; minId++ {
		minIdStr := strconv.Itoa(minId)
		rec = queryRecord(`where id=` + minIdStr)
		if 1 == len(rec) {
			if _last == rec[0].Content {
				fmt.Println("remove" + minIdStr)
				SqliteSql.ExecSqlString(_database, `delete from record where id=`+minIdStr)
			} else if rec[0].Content == "" {
				fmt.Println("remove" + minIdStr)
				SqliteSql.ExecSqlString(_database, `delete from record where id=`+minIdStr)
			} else {
				_last = rec[0].Content
			}
		}
	}
}
