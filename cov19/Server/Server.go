package Server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"public.sunibas.cn/go_public/cov19/Config"
	sqlsint "public.sunibas.cn/go_public/cov19/sql"
	"public.sunibas.cn/go_public/utils/MyServer"
	"time"
)

var sqlsInt sqlsint.SqlsInt

func Server(sc Config.ServerConfig,_sqlsInt sqlsint.SqlsInt)  {
	sqlsInt = _sqlsInt
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer,request,"/public",http.StatusFound)
	})
	InitAction()
	srv := &http.Server{
		Addr:           sc.Port,
		Handler:        nil,
		ReadTimeout:    time.Duration(5) * time.Minute,
		WriteTimeout:   time.Duration(5) * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	err := srv.ListenAndServe()
	defer srv.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func InitAction()  {
	http.HandleFunc("/api", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		writer.Header().Set("content-type", "application/json")             //返回数据格式是json
		params,err := MyServer.PostParams(request)
		if err != nil {
			rObj := MyServer.ReturnObj{}
			rObj.SetFail("解析参数发生错误")
			rObj.PostReturn(writer)
		} else {
			switch params.Method {
			case "selectDI":
				selectDI(writer,params.Content)
				break
			case "insert":
				insert(writer,params.Content)
				break
			case "select":
				selectW(writer,params.Content)
				break
			case "getNames":
				getNames(writer)
				break;
			default:
				rObj := MyServer.ReturnObj{}
				rObj.SetFail("找不到指定方法")
				rObj.PostReturn(writer)
			}
		}
	})
}

type param struct {
	Itype string `json:"itype"`
	From string `json:"from"`
	To string `json:"to"`
	Num string `json:"num"`
	Date string `json:"date"`
	Whos string `json:"whos"`
	Note string `json:"note"`
}
// 解析参数，如果成功则调用 scb 函数，失败则直接对 obj 进行错误返回
func getParam(content string,scb func(param2 param),fcb func(err error)) {
	var p param
	err := json.Unmarshal([]byte(content),&p)
	if err != nil {
		fcb(err)
	} else {
		scb(p)
	}
}

func selectDI(writer http.ResponseWriter, content string) {
	rObj := MyServer.ReturnObj{}
	getParam(content, func(param param) {
		ret := sqlsInt.SelectDI(param.From,param.Date,param.Itype)
		rObj.SetContent(ret)
		rObj.SetSuccess("")
		rObj.PostReturn(writer)
	}, func(err error) {
		rObj.SetFail("参数解析错误" + err.Error())
		rObj.PostReturn(writer)
	})
}
type whereParam struct {
	Table string `json:"table"`
	Where string `json:"where"`
}
func selectW(writer http.ResponseWriter, content string) {
	var p whereParam
	rObj := MyServer.ReturnObj{}
	err := json.Unmarshal([]byte(content),&p)
	if err != nil {
		rObj.SetFail("参数解析错误" + err.Error())
		rObj.PostReturn(writer)
	} else {
		fmt.Println(p.Where)
		ret := sqlsInt.Select(p.Where,p.Table)
		rObj.SetContent(ret)
		rObj.SetSuccess("")
		rObj.PostReturn(writer)
	}
}

func insert(writer http.ResponseWriter, content string) {
	rObj := MyServer.ReturnObj{}
	getParam(content, func(param param) {
		sqlsInt.Insert(param.From,param.To,param.Date,param.Num,param.Itype,param.Whos,param.Note)
		rObj.SetSuccess("插入成功")
		rObj.PostReturn(writer)
	}, func(err error) {
		rObj.SetFail("参数解析错误" + err.Error())
		rObj.PostReturn(writer)
	})
}

func getNames(writer http.ResponseWriter)  {
	rObj := MyServer.ReturnObj{}
	rObj.SetContent(sqlsInt.SelectNames())
	rObj.PostReturn(writer)
}



