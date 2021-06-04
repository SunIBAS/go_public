package main

import (
	"fmt"
	"github.com/alphaqiu/mnemonic"
	"os"
	"public.sunibas.cn/go_public/bit/utils"
	"public.sunibas.cn/go_public/utils/Datas"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"strconv"
	"strings"
	"time"
)

var target = []string{
	"0xea674fdde714fd979de3edf0f56aa9716b898ec8",
	"0x416299aade6443e6f6e8ab67126e65a7f606eef5",
	"0x77ab999d1e9f152156b4411e1f3e2a42dab8cd6d",
	"0xc88f7666330b4b511358b7742dc2a3234710e7b1",
}

var counter int64

func baoliCrash(wlen int, lang mnemonic.Language) {
	var RandCreate func() int64
	var ToWord func(rnd *utils.CreateRnd) string
	RandCreate = utils.DefaultRandCreate
	ToWord = utils.DefaultToWord

	//if rand.Intn(10) < 5 {
	//	RandCreate = utils.DefaultRandCreate
	//	ToWord = utils.DefaultToWord
	//} else {
	//	RandCreate = utils.OtherRandCreate
	//	ToWord = utils.OtherToWord
	//}
	cr := utils.CreateRnd{
		MnwIndex:   make([]int64, wlen),
		Wordlen:    wlen,
		Lang:       lang,
		MNW:        "",
		RandCreate: RandCreate,
		ToWord:     ToWord,
	}
	var addr string
	for {
		cr.GetRnd()
		addr = cr.ToAddress()
		for _, w := range target {
			if addr == w {
				fmt.Println(cr.MNW)
				t := time.Now().Unix()
				DirAndFile.WriteWithWriteString(strconv.FormatInt(t, 10)+".crash.txt", cr.Index2Str()+"\r\n"+cr.MNW)
			} /* else {
				t := time.Now().Unix()
				DirAndFile.WriteWithWriteString(strconv.FormatInt(t,10) + "current.log.txt",cr.MNW)
			}*/
		}
		counter++
	}
}

func main() {
	threadLen := Datas.AtoI(os.Args[1])
	target = DirAndFile.ReadAsFileAsLine(os.Args[2])
	for ind, _ := range target {
		target[ind] = strings.ToLower(target[ind])
		fmt.Println("[" + strconv.Itoa(ind) + "] => " + "[" + target[ind] + "]")
	}
	for i := 0; i < threadLen; i++ {
		go baoliCrash(24, mnemonic.English)
		go baoliCrash(12, mnemonic.English)
		go baoliCrash(24, mnemonic.ChineseSimplified)
		go baoliCrash(12, mnemonic.ChineseSimplified)
		go baoliCrash(24, mnemonic.Japanese)
		go baoliCrash(12, mnemonic.Japanese)
	}
	for {
		time.Sleep(time.Second * 2)
		fmt.Print("\r ====> " + strconv.FormatInt(counter/2, 10) + "Ls/s <====")
		counter = 0
	}
	fmt.Scanln()
}
