package Console

import (
	"bufio"
	"fmt"
	"os"
	"public.sunibas.cn/go_public/utils/Datas"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"public.sunibas.cn/go_public/utils/MyRegex"
	"strconv"
	"strings"
)

func Scanf(a *string) {
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	*a = string(data)
}

func ScanfWhite(line *string) {
	ColorPrint("", FontColor.White)
	//fmt.Scanln(&line)
	Scanf(line)
	ColorPrint("", FontColor.Light_gray)
}

// 新建文件或文件夹
func InputNewPath(input *string, title string, errMsg string) {
	var ptype int
	var err error
	if title == "" {
		title = "请输入文件位置"
	}
	if errMsg == "" {
		errMsg = "指定位置已经存在文件(夹)，请重新输入"
	}
	for {
		fmt.Printf("%s\n请输入>>", title)
		ScanfWhite(input)
		if ptype, err = DirAndFile.PathExistsAndType(*input); err != nil {
			fmt.Println(err)
			return
		} else if ptype != 0 {
			fmt.Println(errMsg)
		} else {
			break
		}
	}
}

// 输入数字
func InputNumber(from, to int, tip string) int {
	var input string
	var ret int
	for {
		InputLine(&input, strconv.Itoa(from), tip)
		if MyRegex.TestIntNumber(input) {
			ret = Datas.AtoI(input)
			if ret >= from && ret <= to {
				break
			}
		} else {
			fmt.Println("输入内容为：", input, "不符合要求，不是数值")
			fmt.Println("要求输入内容为数字，并且在", strconv.Itoa(from), "~", strconv.Itoa(to), "范围中")
		}
	}
	return ret
}

// 获取已有的文件(夹)路径
func InputExistPath(input *string, isFile bool, title string, errMsg string) {
	var ptype int
	var err error
	if title == "" {
		title = "请输入文件位置"
	}
	if errMsg == "" {
		errMsg = "指定位置不存在，请重新输入"
	}
	checkType := -1
	if isFile {
		checkType = -1
	} else {
		checkType = 1
	}
	for {
		fmt.Printf("%s\n请输入>>", title)
		ScanfWhite(input)
		if ptype, err = DirAndFile.PathExistsAndType(*input); err != nil {
			fmt.Println(err)
			return
		} else if ptype != checkType {
			fmt.Println(errMsg)
		} else {
			break
		}
	}
}

// inputRange = [1,2,3]
// inputTip = "请输入 1 2 3 中的一个"
// errMsg = "只能输入 1 2 3"
func InputLineRange(input *string, inputRange []string, inputTip string, errMsg string) {
	get := false
	all := strings.Join(inputRange, "\n")
	for {
		fmt.Println(all)
		fmt.Println(inputTip)
		fmt.Printf("请输入>>")
		ScanfWhite(input)
		for _, n := range inputRange {
			if *input == n {
				get = true
				break
			}
		}
		if get {
			break
		} else {
			fmt.Println(errMsg)
		}
	}
}

/*
从若干个选项中选择若干个（不少于一个）
*/
func SelectRange(inputRange []string, inputTip string, errMsg string) []string {
	selected := []bool{}
	selectAll := false
	selectItem := []string{}
	for i := 0; i < len(inputRange); i++ {
		selected = append(selected, false)
	}
	for {
		fmt.Println(inputTip)
		fmt.Println("all 表示全选")
		fmt.Println("是否选中 项数 项内容")
		fmt.Println(selectAll, "\t0\tall")
		for ind, ir := range inputRange {
			fmt.Println(selected[ind], "\t", ind+1, "\t", ir)
		}
		fmt.Println("输入 -1 表示选择完成")
		selectNum := InputNumber(-1, len(inputRange)+1, "请输入【项数】")
		if selectNum == -1 {
			for ind, s := range selected {
				if s {
					selectItem = append(selectItem, inputRange[ind])
				}
			}
			if len(selectItem) == 0 {
				fmt.Println("至少要选择一个")
			} else {
				break
			}
		}
		if selectNum == 0 {
			if selectAll {
				for i := 0; i < len(inputRange); i++ {
					selected[i] = false
				}
				selectAll = false
			} else {
				for i := 0; i < len(inputRange); i++ {
					selected[i] = true
				}
				selectAll = true
			}
		} else {
			selected[selectNum-1] = !selected[selectNum-1]
		}
	}
	return selectItem
}

func GetInputLineRange(inputRange []string, inputTip string, errMsg string) string {
	var line string
	InputLineRange(&line, inputRange, inputTip, errMsg)
	return line
}

func InputLine(input *string, defaultInput string, inputTip string) {
	fmt.Println(inputTip)
	fmt.Println("回车默认为:" + defaultInput)
	fmt.Printf("请输入>>")
	ScanfWhite(input)
	if len(*input) == 0 {
		*input = defaultInput
	}
}

var line string

func InputToContinue() {
	fmt.Println("输入回车继续")
	Scanf(&line)
}
