package main

import (
	"fmt"
	"os"
	"path"
	"public.sunibas.cn/go_public/utils/Console"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"strings"
)

func main() {
	fmt.Println("请确保安装了 GDAL 并配置到环境变量(PATH)中")
	fmt.Println("GDAL 下载地址：http://download.osgeo.org/gdal/")
	var linesDir string
	Console.Clear()
	fmt.Println("样本点文件夹一般包含多个样本，每个样本在一个文件夹中，格式可能如下")
	fmt.Println("c:")
	fmt.Println("|- 1 <- 这个就是要输入的文件夹")
	fmt.Println("|---- 0000001100_V1")
	fmt.Println("|---- ---- 0000001100.tif")
	fmt.Println("|---- ---- 0000001100_V1_LINE.shp")
	fmt.Println("|---- ---- 0000001100_V1_LINE.xxx")
	fmt.Println("|---- 0000001101_V1")
	fmt.Println("|---- ---- 0000001101.tif")
	fmt.Println("|---- ---- 0000001101_V1_LINE.shp")
	fmt.Println("|---- ---- 0000001101_V1_LINE.xxx")
	fmt.Println("|---- ...")
	Console.InputLine(&linesDir, "c:\\1", "请输入样本点的文件夹（全路径）")
	var outputDir string
	Console.InputLine(&outputDir, linesDir+"_tif", "请输入保存样本点面化文件的文件夹（最好为空）")
	var burn string
	Console.InputLine(&burn, "1", "gdal_rasterize 的 burn 参数，种植结果默认为 1")
	var txt_file = path.Join(outputDir, "train_line.txt")

	fileAndDirs := DirAndFile.GetSubDirOrFile(linesDir)
	//shpFiles := []string{}
	//tifFiles := []string{}
	shpFile := ""
	srcTifFile := ""
	tifFile := ""
	cmdContent := []string{"chcp 65001"}
	txtContent := []string{}
	for _, fd := range fileAndDirs {
		if !fd.File {
			sfds := DirAndFile.GetSubDirOrFile(fd.FullPath)
			for _, sfd := range sfds {
				if sfd.File {
					ext := strings.ToLower(sfd.Name[len(sfd.Name)-4:])
					if ext == ".shp" {
						shpFile = sfd.FullPath
						cmdContent = append(cmdContent, "gdal_rasterize.exe -burn "+burn+" -ts 1000 1000 -init 0 -ot Byte \""+shpFile+"\" \""+tifFile+"\"")
					} else if ext == ".tif" {
						srcTifFile = sfd.FullPath
						tifFile = path.Join(outputDir, sfd.Name[:len(sfd.Name)-4]+"_LINE.tif")
						txtContent = append(txtContent, srcTifFile+" "+tifFile)
					}
				}
			}
		}
	}
	cmdContent = append(cmdContent, "exit")
	if t, err := DirAndFile.PathExistsAndType(outputDir); err == nil {
		if t == DirAndFile.File {
			panic("存在同名文件，无法创建文件夹")
		} else if t == DirAndFile.NoExist {
			os.Mkdir(outputDir, os.ModeDir)
		}
	} else {
		panic(err)
	}
	DirAndFile.WriteWithIOUtil(txt_file, strings.Join(txtContent, "\r\n"))
	DirAndFile.WriteWithIOUtil("_tmp_.bat", strings.Join(cmdContent, "\r\n"))
	Console.RunBatFile("_tmp_.bat")
}
