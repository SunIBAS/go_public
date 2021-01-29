package office

import "github.com/360EntSecGroup-Skylar/excelize/v2"

type NewExcelFile struct {
	File *excelize.File
	Path string
}

func (nef *NewExcelFile) Init() *excelize.File {
	nef.File = excelize.NewFile()
	return nef.File
}
func (nef NewExcelFile) Save() {
	if err := nef.File.SaveAs(nef.Path); err != nil {
		panic(err)
	}
}
func (nef NewExcelFile) SaveAndDeleteSheet1() {
	if err := nef.File.SaveAs(nef.Path); err != nil {
		panic(err)
	}
	// 删除 sheet1
	ofile, err := excelize.OpenFile(nef.Path)
	if err != nil {
		panic(err)
	}
	ofile.DeleteSheet("Sheet1")
	if err := ofile.SaveAs(nef.Path); err != nil {
		println(err.Error())
	}
}
