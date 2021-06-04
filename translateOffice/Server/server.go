package Server

import (
	"fmt"
	"net/http"
	"public.sunibas.cn/go_public/translateOffice/db"
	"public.sunibas.cn/go_public/utils/Console"
	"public.sunibas.cn/go_public/utils/Datas"
	"public.sunibas.cn/go_public/utils/MyRegex"
	"strconv"
	"strings"
	"time"
)

var database *db.LangDb

func StartServer(db_ *db.LangDb, port string) {
	database = db_
	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        nil,
		ReadTimeout:    time.Duration(5) * time.Minute,
		WriteTimeout:   time.Duration(5) * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//Utils.ReturnHTML(writer,Utils.TranslateFileHTML)
		content := "<html><body><script>window.onload=function() {location.href='/1';}</script></body></html>"
		ok := false
		offset := 0
		if request.RequestURI != "/" {
			p := request.RequestURI[1:]
			if len(p) > 0 {
				if MyRegex.TestNumber(p) {
					offset = Datas.AtoI(p) - 1
					ok = true
				}
			}
		}
		if ok {
			langs := database.Select("", "limit 1000 OFFSET "+strconv.Itoa(offset*1000))
			line := ""
			if len(langs) > 0 {
				for _, lang := range langs {
					line += `<div id="` + lang.Id + `">` + lang.Content + `</div>`
				}
				content = strings.Replace(TranslateHTML, "$content$", line, 1)
			} else {
				content = "<html><head><meta charset=\"UTF-8\"></head><body>" +
					"nothing more" +
					"</body></html>"
			}
		}
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		writer.Header().Set("Content-Length", strconv.Itoa(len(content)))
		fmt.Fprint(writer, content)
	})
	Console.RunBatFile("http://localhost:" + port)
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
