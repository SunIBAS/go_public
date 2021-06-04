package InsertJson

import (
	"encoding/json"
	"fmt"
	"public.sunibas.cn/go_public/translateOffice/db"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"strconv"
)

type rec struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}
type t struct {
	Rec []rec `json:"rec"`
}

func InsertJson(db_ *db.LangDb, jsonFile string) {
	tt := t{}
	//content := "{\"rec\":" + DirAndFile.ReadAll(jsonFile) + "}"
	err := json.Unmarshal([]byte("{\"rec\":"+DirAndFile.ReadAll(jsonFile)+"}"), &tt)
	ids := []string{}
	contents := []string{}
	if err == nil {
		for _, r := range tt.Rec {
			ids = append(ids, r.Id)
			contents = append(contents, r.Content)
			if len(ids) > 300 {
				db_.UpdateTrans(ids, contents)
				ids = []string{}
				contents = []string{}
			}
		}
		if len(ids) > 0 {
			db_.UpdateTrans(ids, contents)
		}
		fmt.Println("成功插入" + strconv.Itoa(len(tt.Rec)) + "条记录")
	} else {
		panic(err)
	}
}
