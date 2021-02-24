package Datas

type MultiTimeTestContain struct {
	items map[interface{}]int
}

func (mttc *MultiTimeTestContain) AddStringItem(items []string) {
	if mttc.items == nil {
		mttc.items = map[interface{}]int{}
	}
	for _, i := range items {
		mttc.items[i] = 1
	}
}

func (mttc MultiTimeTestContain) TestString(item string) bool {
	return mttc.items[item] != 0
}

func ArrayContainItemString(arr []string, item string) bool {
	contain := false
	for _, i := range arr {
		if i == item {
			contain = true
		}
	}
	return contain
}
