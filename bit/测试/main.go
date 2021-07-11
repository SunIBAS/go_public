package main

import (
	"crypto/rand"
	"fmt"
	"github.com/alphaqiu/mnemonic"
	"github.com/alphaqiu/mnemonic/entropy"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"log"
	"public.sunibas.cn/go_public/bit/utils"
	"strconv"
	"strings"
	"time"
)

var (
	continue__ = false
	runTime    = time.Millisecond * 100
	sleepTime  = time.Millisecond * 300
)

func setTimeout() {
	rt := time.NewTicker(runTime)
	times := -1
	for {
		select {
		case <-rt.C:
			times++
			times %= 3
			if times != 0 {
				continue__ = false
				fmt.Println("开始休息")
			} else {
				// 开 1 分钟关 一个钟
				continue__ = true
				fmt.Println("休息时间结束")
			}
			// rt.Stop()
		}
	}
}

func main() {
	rndMneProcess()
}

func testCreateRnd() {
	cr := utils.CreateRnd{
		MnwIndex: make([]int64, 24),
		Wordlen:  24,
		Lang:     mnemonic.English,
		MNW:      "",
	}
	fmt.Println(cr.GetRnd())
	fmt.Println(cr.ToAddress())
}

// 理清楚随机助记词的产生和编码位地址的过程
func rndMneProcess() {
	//
	seed := bip39.NewSeed("cake apple borrow silk endorse fitness top denial coil riot stay wolf luggage oxygen faint major edit measure invite love trap field dilemma oblige", "")
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		log.Fatal(err)
	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0") //最后一位是同一个助记词的地址id，从0开始，相同助记词可以生产无限个地址
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}
	address := account.Address.Hex()
	//privateKey, _ := wallet.PrivateKeyHex(account)
	//publicKey, _ := wallet.PublicKeyHex(account)
	fmt.Println(strings.ToLower(address))
}

func rndMne() {
	// 128  |   4      | 132           | 12    |
	// -----|----------|---------------|-------|
	// 160  |   5      | 165           | 15    |
	// -----|----------|---------------|-------|
	// 192  |   6      | 198           | 18    |
	// -----|----------|---------------|-------|
	// 224  |   7      | 231           | 21    |
	// -----|----------|---------------|-------|
	// 256  |   8      | 264           | 24    |
	lens := []int{128, 160, 192, 224, 256}
	langs := []mnemonic.Language{
		mnemonic.English,
		mnemonic.ChineseSimplified,
		mnemonic.ChineseTraditional,
		mnemonic.French,
		mnemonic.Italian,
		mnemonic.Japanese,
		mnemonic.Korean,
	}
	ws := ""
	for _, l := range lens {
		for _, la := range langs {
			mnw, _ := mnemonic.NewRandom(l, la)
			fmt.Println(mnw.Sentence())
			ws += mnw.Sentence() + "\r\n"
		}
	}
	//DirAndFile.WriteWithIOUtil("D:\\codes\\node\\EarnMoneyWay\\kengdie\\mns.txt",ws)
}

func testSetTimeout() {
	go setTimeout()
	cs := false
	ds := false
	for {
		if continue__ {
			ds = false
			if !cs {
				cs = true
				fmt.Print("c")
			}
		} else {
			cs = false
			if !ds {
				ds = true
				fmt.Print("d")
			}
		}
	}
}

func testRandRead() {
	bitSize := 32
	bytes := make([]byte, bitSize/8)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	fmt.Println(bytes)
	bits := entropy.CheckSummed(bytes)
	fmt.Println(bits)
}

func checkkMnemonic() {
	v, _ := mnemonic.IsMnemonicValid(mnemonic.English, "ability able about above absent absorb abstract absurd abuse afford second")
	fmt.Println(v)
}

func timeStamp() {
	t := time.Now().Unix()
	fmt.Println(strconv.FormatInt(t, 10))
}

func int64Arr2str() {
	var currentIndex = []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 38, 1127}
	var strs = []string{}
	for i := 0; i < len(currentIndex); i++ {
		strs = append(strs, strconv.FormatInt(currentIndex[i], 10))
	}
	fmt.Println(strings.Join(strs, ","))
}
