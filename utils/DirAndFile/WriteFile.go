package DirAndFile

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func WriteWithIOUtilByte(name string, data []byte) {
	if ioutil.WriteFile(name, data, 0644) == nil {
		//fmt.Println("写入文件成功:")
	}
}
func WriteWithIOUtil(name, content string) {
	data := []byte(content)
	if ioutil.WriteFile(name, data, 0644) == nil {
		//fmt.Println("写入文件成功:")
	}
}

func WriteWithWriteStringCsv(filename, wireteString string) {
	WriteWithWriteString(filename, "\xEF\xBB\xBF"+wireteString)
}
func WriteWithWriteString(filename, wireteString string) {
	var f *os.File
	var err1 error
	f, err1 = os.Create(filename) //创建文件
	defer f.Close()
	checkErr(err1)
	n, err1 := io.WriteString(f, wireteString) //写入文件(字符串)
	checkErr(err1)
	fmt.Printf("写入 %d 个字节n", n)
}

func WriteWithBufio(filename string, wireteString string) {
	var f *os.File
	var err1 error
	f, err1 = os.Create(filename) //创建文件
	checkErr(err1)
	w := bufio.NewWriter(f) //创建新的 Writer 对象
	n4, err1 := w.WriteString(wireteString)
	checkErr(err1)
	fmt.Printf("写入 %d 个字节n", n4)
	w.Flush()
	f.Close()
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
