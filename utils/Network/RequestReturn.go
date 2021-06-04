package Network

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ReturnObj struct {
	Code    int
	Message string
	Content string
}

func (rObj *ReturnObj) SetStringContent(obj string) *ReturnObj {
	rObj.Content = obj
	return rObj
}
func (rObj *ReturnObj) SetContent(obj interface{}) *ReturnObj {
	jsonStr, _ := json.MarshalIndent(obj, "", "")
	rObj.Content = string(jsonStr)
	return rObj
}

func (rObj *ReturnObj) SetError(msg string) *ReturnObj {
	rObj.Code = Error
	if len(msg) == 0 {
		rObj.Message = "fail"
	} else {
		rObj.Message = msg
	}
	return rObj
}
func (rObj *ReturnObj) SetFail(msg string) *ReturnObj {
	rObj.Code = Fail
	if len(msg) == 0 {
		rObj.Message = "fail"
	} else {
		rObj.Message = msg
	}
	return rObj
}
func (rObj *ReturnObj) SetSuccess(msg string) *ReturnObj {
	rObj.Code = Success
	if len(msg) == 0 {
		rObj.Message = "success"
	} else {
		rObj.Message = msg
	}
	return rObj
}

func (rObj *ReturnObj) Send(w http.ResponseWriter) *ReturnObj {
	jsonStr, _ := json.MarshalIndent(rObj, "", "")
	content := string(jsonStr)
	content = strings.Replace(content, "\\u003c", "<", -1)
	content = strings.Replace(content, "\\u003e", ">", -1)
	content = strings.Replace(content, "\\u0026", "&", -1)
	fmt.Fprintf(w, content)
	return rObj
}

func PostReturn(w http.ResponseWriter, rObj ReturnObj) {
	jsonStr, _ := json.MarshalIndent(rObj, "", "")
	fmt.Fprintf(w, string(jsonStr))
}

func DownloadFile(w http.ResponseWriter, filename string) {
	file, _ := os.Open(filename)
	defer file.Close()
	fileHeader := make([]byte, 512)
	file.Read(fileHeader)
	fileStat, _ := file.Stat()
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	file.Seek(0, 0)
	io.Copy(w, file)
}
func DownloadEmptyFile(w http.ResponseWriter, filename, reason string) {
	content := "file [" + filename + "] is not exist, reason is [" + reason + "]"
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	fmt.Fprint(w, content)
}

func DownloadContent(w http.ResponseWriter, filename, content string) {
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	fmt.Fprint(w, content)
}

func ReturnHTML(w http.ResponseWriter, content string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	fmt.Fprint(w, content)
}
