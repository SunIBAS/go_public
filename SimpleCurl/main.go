package main

import (
	"compress/gzip"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"strings"
)

var (
	method     string
	url        string
	browser    bool
	header     map[string]string
	parameters string
	outputFile string
)

func usage() {
	fmt.Println("Usage:")
	fmt.Println("scurl [method=get] [browser] [h=k;v] [p=] [o=] url")
	fmt.Println("\tmethod 可选有 get post")
	fmt.Println("\tbrowser 表示是否追加模拟浏览器头部USER-AGENT")
	fmt.Println("\th header 信息，key-value 用 ; 隔开")
	fmt.Println("\tp 请求参数，例如 p=name=ibas&age=12")
	fmt.Println("\to 保存请求内容到文件，例如 o=a.txt")
	fmt.Println("\turl 请求链接")
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func main() {
	url = "https://sp0.baidu.com/8aQDcjqpAAV3otqbppnN2DJv/api.php?query=2001%E5%B9%B41%E6%9C%88&co=&resource_id=39043&t=1625036257897&ie=utf8&oe=gbk&cb=op_aladdin_callback&format=json&tn=wisetpl&cb=callback&_=1624981279114"
	method = "get"
	browser = false
	header = map[string]string{}
	parameters = ""
	outputFile = ""
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			if len(os.Args[i]) > len("method=") &&
				os.Args[i][:len("method=")] == "method=" {
				method = os.Args[i][len("method="):]
			} else if os.Args[i] == "browser" {
				browser = true
			} else if len(os.Args[i]) > len("http://") && os.Args[i][:len("http")] == "http" {
				url = os.Args[i]
			} else if len(os.Args[i]) > len("h=") && os.Args[i][:len("h=")] == "h=" {
				sp := strings.Split(os.Args[i][len("h="):], ";")
				if len(sp) == 2 {
					header[sp[0]] = sp[1]
				}
			} else if len(os.Args[i]) > len("p=") && os.Args[i][:len("p=")] == "p=" {
				parameters = os.Args[i][len("p="):]
			} else if len(os.Args[i]) > len("o=") && os.Args[i][:len("o=")] == "o=" {
				outputFile = os.Args[i][len("o="):]
			}
		}
		if method != "get" && method != "post" {
			method = "get"
		}
		method = strings.ToUpper(method)
		toReq()
	} else {
		usage()
	}
}

func toReq() {
	//传递的参数: qq=908204694
	dataReader := strings.NewReader(parameters)
	newReq, err := http.NewRequest(method, url, dataReader)
	if err != nil {
		return
	}
	if browser {
		newReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
	}
	//newReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range header {
		newReq.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(newReq)

	if err != nil {
		fmt.Println(err)
	}
	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		fmt.Println("gzip")
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return
		}
	} else {
		reader = resp.Body
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(reader)
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	var str string
	if strings.Contains(contentType, "charset=gbk") {
		str = ConvertToString(string(body), "gbk", "utf-8")
	} else {
		str = string(body)
	}
	if err != nil {
		// handle error
	}
	fmt.Println("outputFile = " + outputFile)
	if len(outputFile) > 0 {
		DirAndFile.WriteWithIOUtil(outputFile, str)
	} else {
		//打印返回值
		fmt.Println(str)
	}
}
