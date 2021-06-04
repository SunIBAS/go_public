package Datas

import (
	"strconv"
	"strings"
)

// 将 16 进制的 uuid 转化为 36 进制的 uuid
// ffffffff-ffff-ffff-ffff-ffffffffffff
// 1z141z3-2rrvthnxtr-2rrvthnxtr
func UUID16236(s string) string {
	as := ""
	i64, _ := strconv.ParseInt(s[0:8], 16, 64)
	as += "-" + strconv.FormatInt(i64, 36)
	i64, _ = strconv.ParseInt(strings.Replace(s[9:23], "-", "", 2), 16, 64)
	as += "-" + strconv.FormatInt(i64, 36)
	i64, _ = strconv.ParseInt(s[24:], 16, 64)
	as += "-" + strconv.FormatInt(i64, 36)
	return as[1:]
}

// 将 36 进制的 uuid 转化为 16 进制的 uuid
// ffffffff-ffff-ffff-ffff-ffffffffffff
// 1z141z3-2rrvthnxtr-2rrvthnxtr
func UUID36216(s string) string {
	blen := []int{
		8, 12, 12,
	}
	zeros := "000000000000"
	as := []string{}
	ss := strings.Split(s, "-")
	for _, is := range ss {
		i64, _ := strconv.ParseInt(is, 36, 64)
		as = append(as, strconv.FormatInt(i64, 16))
	}
	for ind, _ := range blen {
		as[ind] = zeros[len(as[ind]):blen[ind]] + as[ind]
	}
	return as[0] + "-" + as[1][0:4] + "-" + as[1][4:8] + "-" + as[1][8:12] + "-" + as[2]
}
