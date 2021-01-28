package MyServer

import (
	"encoding/json"
	"net/http"
)

type Parmas struct {
	Method string
	Content string
}

func PostParams(r * http.Request) (Parmas,error) {
	var p Parmas
	err := json.NewDecoder(r.Body).Decode(&p)
	return p,err
}
