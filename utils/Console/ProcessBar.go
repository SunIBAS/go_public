package Console

import (
	"fmt"
	"strconv"
	"strings"
)

// https://www.zhihu.com/question/29738147
func Bar(vl int, width int) string {
	return fmt.Sprintf("%s%*c", strings.Repeat("█", vl/10), vl/10-width+1,
		([]rune(" ▏▎▍▌▋▋▊▉█"))[vl%10])
}

type EqBar struct {
	Total   int // 总的量
	Percent int // 当前的量
	Size    int // 进度条分为几部分
}

// 显示并加一
func (eb *EqBar) ShowPlusOne() {
	eb.Percent++
	eb.Show()
}

// 显示（这个是以百分数输出的）
func (eb EqBar) Show() {
	str := "[" + eqBar(eb.Size*eb.Percent/eb.Total, eb.Size) + "] " +
		strconv.Itoa(100*eb.Percent/eb.Total) + "%"
	fmt.Printf("\r%s", str)
}
func (eb EqBar) GetShow() string {
	return "[" + eqBar(eb.Size*eb.Percent/eb.Total, eb.Size) + "] " +
		strconv.Itoa(100*eb.Percent/eb.Total) + "%"
}

// 这个也是显示并加一，但是是调用下面这个显示
func (eb *EqBar) ShowNumPlusOne() {
	eb.Percent++
	eb.ShowNum()
}

// 以分数形式显示
func (eb EqBar) ShowNum() {
	cur := eb.Size * eb.Percent / eb.Total
	str := "[" + eqBar(cur, eb.Size) + "] " +
		strconv.Itoa(eb.Percent) + "/" + strconv.Itoa(eb.Total)
	fmt.Printf("\r%s", str)
}
func (eb EqBar) GetShowNum() string {
	cur := eb.Size * eb.Percent / eb.Total
	return "[" + eqBar(cur, eb.Size) + "] " +
		strconv.Itoa(eb.Percent) + "/" + strconv.Itoa(eb.Total)
}
func eqBar(count, size int) string {
	str := ""
	for i := 0; i < size; i++ {
		if i < count {
			str += "="
		} else {
			str += " "
		}
	}
	return str
}

type TwoEqBar struct {
	Total     int // 总的量
	Percent   int // 当前的量
	Size      int // 进度条分为几部分
	Asize     int // 第一部分长度
	Bsize     int // 第二部分长度
	APartChar string
	BPartChar string
	OtherChar string
}

// 显示并加一
func (eb *TwoEqBar) ShowPlusOne() {
	eb.Percent++
	eb.Show()
}

// 显示（这个是以百分数输出的）
func (eb TwoEqBar) Show() {
	str := "[" + eb.twoEqBar() + "] " +
		strconv.Itoa(100*eb.Percent/eb.Total) + "%"
	fmt.Printf("\r%s", str)
}
func (eb TwoEqBar) GetShow() string {
	return "[" + eb.twoEqBar() + "] " +
		strconv.Itoa(100*eb.Percent/eb.Total) + "%"
}

// 这个也是显示并加一，但是是调用下面这个显示
func (eb *TwoEqBar) ShowNumPlusOne() {
	eb.Percent++
	eb.ShowNum()
}

// 以分数形式显示
func (eb TwoEqBar) ShowNum() {
	str := "[" + eb.twoEqBar() + "] " +
		strconv.Itoa(eb.Percent) + "/" + strconv.Itoa(eb.Total)
	fmt.Printf("\r%s", str)
}
func (eb TwoEqBar) GetShowNum() string {
	return "[" + eb.twoEqBar() + "] " +
		strconv.Itoa(eb.Percent) + "/" + strconv.Itoa(eb.Total)
}
func (eb TwoEqBar) TEBString() string {
	return eb.twoEqBar()
}
func (eb TwoEqBar) twoEqBar() string {
	count := eb.Size * eb.Percent / eb.Total
	str := ""
	for i := 0; i < eb.Size; i++ {
		if i < count {
			if i < eb.Asize {
				str += eb.APartChar
			} else {
				str += eb.BPartChar
			}
		} else {
			str += eb.OtherChar
		}
	}
	return str
}
