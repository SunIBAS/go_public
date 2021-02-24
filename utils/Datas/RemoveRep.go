package Datas

func RemoveRepByMap(slc []interface{}) []interface{} {
	result := []interface{}{}
	tempMap := map[interface{}]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

func RemoveRepByMapString(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

func StatisticsMapString(slc []string) map[string]int {
	ret := map[string]int{}
	for _, s := range slc {
		if _, ok := ret[s]; ok {
			ret[s]++
		} else {
			ret[s] = 1
		}
	}
	return ret
}
