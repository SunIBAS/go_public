package DirAndFile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func GetDirSize(dirPath string) int64 {
	var dirSize int64 = 0
	flist, e := ioutil.ReadDir(dirPath)
	if e != nil {
		fmt.Println("read ", dirPath)
		fmt.Println("read file error")
		return 0
	}
	for _, f := range flist {
		if f.IsDir() {
			dirSize = GetDirSize(dirPath+"/"+f.Name()) + dirSize
		} else {
			dirSize = f.Size() + dirSize
		}
	}
	return dirSize
}

type FileNode struct {
	Name     string
	File     bool
	FullPath string
	// 相对路径
	RPath string
}

func GetSubDirOrFile(dirPath string) []FileNode {
	fileNode := []FileNode{}
	flist, e := ioutil.ReadDir(dirPath)
	if e != nil {
		fmt.Println("read ", dirPath)
		fmt.Println("read file error")
		return fileNode
	}
	for _, f := range flist {
		fileNode = append(fileNode, FileNode{
			Name:     f.Name(),
			File:     !f.IsDir(),
			FullPath: path.Join(dirPath, f.Name()),
		})
	}
	return fileNode
}

// 路径分隔符，windows下应该是 \
const Filepath_Separator = string(rune(filepath.Separator))

func FilePathWalkDir(root string) ([]FileNode, error) {
	fileNode := []FileNode{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path == root {
			return nil
		}
		fileNode = append(fileNode, FileNode{
			Name:     info.Name(),
			File:     !info.IsDir(),
			RPath:    path[len(root):],
			FullPath: path,
		})
		return nil
	})
	return fileNode, err
}
