package DirAndFile

import "os"

// 只能判断一个路径是否存在 如果 true 表示 不存在，false 表示 存在
func PathNotExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return false
	}
	if os.IsNotExist(err) {
		// 只有这里能确定文件不存在
		return true
	}
	return false
}

// 如果返回的错误为nil,说明文件或文件夹存在
// 如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
// 如果返回的错误为其它类型,则不确定是否在存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

const (
	File    = -1
	NoExist = 0
	Dir     = 1
)

// -1 是文件，1 是文件夹，0表示不存在
func PathExistsAndType(path string) (int, error) {
	s, err := os.Stat(path)
	if err == nil {
		if s.IsDir() {
			return 1, nil
		} else {
			return -1, nil
		}
	}
	if os.IsNotExist(err) {
		return 0, nil
	}
	return 0, err
}
