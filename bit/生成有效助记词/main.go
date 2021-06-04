package main

import (
	"fmt"
	"github.com/alphaqiu/mnemonic"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"strconv"
	"strings"
	"time"
)

var currentIndex = []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 92, 321, 1297}
var pointIndex = 11
var totalWords int64 = 2048

func main() {
	okmnw := []string{}
	for {
		if !next() {
			break
		}
		mnw := toWord()
		//fmt.Println(mnw)
		if ok, _ := mnemonic.IsMnemonicValid(mnemonic.English, mnw); ok {
			okmnw = append(okmnw, mnw)
			if len(okmnw) == 200 {
				tmp := okmnw
				okmnw = []string{}
				t := time.Now().Unix()
				//fmt.Println(strconv.FormatInt(t,10))
				DirAndFile.WriteWithIOUtil(
					"D:\\codes\\node\\EarnMoneyWay\\tmp\\getbalances\\ok_zjc\\"+strconv.FormatInt(t, 10)+".go.txt",
					currentIndex2str()+strings.Join(tmp, "\r\n"))
				fmt.Println(currentIndex2str())
				//break
			}
		}
	}
	fmt.Println("over")
}

func currentIndex2str() string {
	var strs = []string{}
	for i := 0; i < len(currentIndex); i++ {
		strs = append(strs, strconv.FormatInt(currentIndex[i], 10))
	}
	return "currentIndex=" + strings.Join(strs, ",") + "\r\n"
}

func plusOne() bool {
	currentIndex[pointIndex]++
	if currentIndex[pointIndex] >= totalWords {
		return false
	}
	for i := 0; i < pointIndex; i++ {
		if currentIndex[i] == currentIndex[pointIndex] {
			return plusOne()
		}
	}
	return true
}

func next() bool {
	if !plusOne() {
		pointIndex--
	}
	for {
		if !plusOne() {
			pointIndex--
			if pointIndex == -1 {
				break
			}
		} else {
			if pointIndex == len(currentIndex)-1 {
				break
			}
			pointIndex++
			currentIndex[pointIndex] = 0
		}
	}
	if pointIndex == -1 {
		return false
	} else {
		return true
	}
}

func toWord() string {
	var ws, _ = mnemonic.GetWord(mnemonic.English, currentIndex[0])

	for i := 1; i < len(currentIndex); i++ {
		w, _ := mnemonic.GetWord(mnemonic.English, currentIndex[i])
		ws += " " + w
	}
	return ws
}
