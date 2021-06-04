package DirAndFile

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func ReadAsFileAsLine(filename string) []string {
	var lines []string
	fmt.Println("read file : ", filename)
	if file, err := os.Open(filename); err != nil {
		panic(err)
	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}
		return lines
	}
}

func ReadAll(filename string) string {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("read fail", err)
	}
	return string(f)
}
