package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port     int    `json:"port"`
	ApiName  string `json:"apiName"`
	Action   string `json:"action"`
	FilePath string `json:"filePath"`
}

func main() {
	config := Config{}
	if len(os.Args) == 2 {
		jsonContnet := strings.Join(DirAndFile.ReadAsFileAsLine(os.Args[1]), "\n")
		if err := json.Unmarshal([]byte(jsonContnet), &config); err != nil {
			panic(err)
		} else {
			createApi(config)
		}
	}
}

func createApi(config Config) {
	srv := &http.Server{
		Addr:           ":" + strconv.Itoa(config.Port),
		Handler:        nil,
		ReadTimeout:    time.Duration(5) * time.Minute,
		WriteTimeout:   time.Duration(5) * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	action := func(a string) {}
	if config.Action == "file" {
		action = func(content string) {
			Action_file(content, config.FilePath)
		}
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")  //允许访问所有域
		writer.Header().Add("Access-Control-Allow-Headers", "*") //header的类型
		writer.Header().Set("content-type", "application/json")  //返回数据格式是json
		s, _ := ioutil.ReadAll(request.Body)
		action(string(s))
		fmt.Fprintf(writer, "ok")
	})
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Action_file(content, filePath string) {
	DirAndFile.WriteWithIOUtil(filePath, content)
}
