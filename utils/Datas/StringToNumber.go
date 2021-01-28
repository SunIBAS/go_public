package Datas

import "strconv"

func AtoI(str string) int {
	if v, e := strconv.Atoi(str); e != nil {
		panic(e)
	} else {
		return v
	}
}

func AtoI64(str string) int64 {
	return int64(AtoI(str))
}

func AtoF(str string) float64 {
	if fl, err := strconv.ParseFloat(str, 64); err != nil {
		panic(err)
	} else {
		return fl
	}
}
