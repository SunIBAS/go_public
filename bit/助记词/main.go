package main

import (
	"fmt"
	"github.com/alphaqiu/mnemonic"
	"net/http"
	"time"
)

func main() {
	createApi()
}

func createApi() {
	srv := &http.Server{
		Addr:           ":8765",
		Handler:        nil,
		ReadTimeout:    time.Duration(5) * time.Minute,
		WriteTimeout:   time.Duration(5) * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/english", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")  //允许访问所有域
		writer.Header().Add("Access-Control-Allow-Headers", "*") //header的类型
		writer.Header().Set("content-type", "application/text")  //返回数据格式是json
		m, _ := mnemonic.NewRandom(256, mnemonic.English)
		fmt.Fprint(writer, m.Sentence())
	})
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
