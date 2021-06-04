package main

import (
	"fmt"
	"github.com/alphaqiu/mnemonic"
	"math/rand"
	"os"
	"path"
	"public.sunibas.cn/go_public/utils/Datas"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"strconv"
	"strings"
	"time"
)

var (
	wordlen   int
	hashlen   int
	savepath  string
	mnwIndex  []int64
	findex    int
	lastIndex int
	fileCount int
)

func main() {
	lastIndex = 0
	if len(os.Args) > 5 {
		wordlen = Datas.AtoI(os.Args[1])
		savepath = os.Args[2]
		if wordlen == 12 {
			hashlen = 128
		} else if wordlen == 15 {
			hashlen = 160
		} else if wordlen == 18 {
			hashlen = 192
		} else if wordlen == 21 {
			hashlen = 224
		} else if wordlen == 24 {
			hashlen = 256
		}
		mnwIndex = make([]int64, wordlen)
		findex = Datas.AtoI(os.Args[4])
		fileCount = Datas.AtoI(os.Args[5])
		if os.Args[3] == "f" {
			doit()
		} else if os.Args[3] == "t" {
			doitRnd()
		}
	} else {
		fmt.Println("program.exe len savepath r")
		fmt.Println("len = [12|15|18|21|24]")
		fmt.Println("savepath 保存文件位置")
		fmt.Println("r = [t|f] t 表示真随机")
		fmt.Println("fn 文件索引")
		fmt.Println("n 文件数量")
	}
}

func getFileName() string {
	t := time.Now().Unix()
	if findex != 0 {
		lastIndex = lastIndex + 1
		lastIndex = lastIndex % findex
		return path.Join(savepath, strconv.FormatInt(t, 10)+"."+strconv.Itoa(lastIndex)+".go.txt")
	} else {
		return path.Join(savepath, strconv.FormatInt(t, 10)+".go.txt")
	}
}

func doit() {
	fmt.Println("rand na")
	okmnw := []string{}
	for {
		mnw, _ := mnemonic.NewRandom(hashlen, mnemonic.English)
		okmnw = append(okmnw, mnw.Sentence())
		if len(okmnw) == 200 {
			tmp := okmnw
			okmnw = []string{}
			//fmt.Println(strconv.FormatInt(t,10))
			DirAndFile.WriteWithIOUtil(
				getFileName(),
				"line 1\r\n"+strings.Join(tmp, "\r\n"))
			time.Sleep(time.Second)
			fileCount--
			if fileCount == 0 {
				break
			}
		}
	}
	fmt.Println("over")
}

func toWord() string {
	var ws, _ = mnemonic.GetWord(mnemonic.English, mnwIndex[0])

	for i := 1; i < len(mnwIndex); i++ {
		w, _ := mnemonic.GetWord(mnemonic.English, mnwIndex[i])
		ws += " " + w
	}
	return ws
}
func getRnd() string {
	ok := false
	currentLen := 0
	for {
		n := rand.Int63n(2048)
		ok = true
		for i := 0; i < currentLen; i++ {
			if mnwIndex[i] == n {
				ok = false
				break
			}
		}
		if ok {
			mnwIndex[currentLen] = n
			currentLen++
		}
		if currentLen == wordlen {
			break
		}
	}
	mnw := toWord()
	if ok, _ := mnemonic.IsMnemonicValid(mnemonic.English, mnw); ok {
		return mnw
	} else {
		return getRnd()
	}
}
func doitRnd() {
	fmt.Println("rand pr")
	okmnw := []string{}
	for {
		mnw := getRnd()
		okmnw = append(okmnw, mnw)
		if len(okmnw) == 200 {
			tmp := okmnw
			okmnw = []string{}
			//fmt.Println(strconv.FormatInt(t,10))
			//fmt.Println("doit")
			DirAndFile.WriteWithIOUtil(
				getFileName(),
				"line 1\r\n"+strings.Join(tmp, "\r\n"))
			time.Sleep(time.Second)
			fileCount--
			if fileCount == 0 {
				break
			}
		}
	}
	fmt.Println("over")
}
