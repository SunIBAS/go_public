package ParseOffice

import (
	"github.com/subchen/go-xmldom"
	"log"
	"os"
	"public.sunibas.cn/go_public/utils/Datas"
	"strings"
)

type RpExcel struct {
	InFile       string
	OutFile      string
	XmlName      string
	SheetXmlName string
	SheetName    []string
	Words        []string
	Xml          string
	SheetXml     string
	Merge        bool
	Logger       *log.Logger
}

// 初始化文档对象
func (re *RpExcel) Init() {
	re.XmlName = "xl/sharedStrings.xml"
	re.SheetXmlName = "xl/workbook.xml"
	re.Xml = ""
	re.Words = []string{}
	re.SheetName = []string{}
	if re.Logger == nil {
		re.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}
}

// 获取文档文字信息
func (re *RpExcel) GetWords() []string {
	getW := func() {
		doc := xmldom.Must(xmldom.ParseXML(re.Xml))
		root := doc.Root
		for _, n := range root.Query("//t") {
			re.Words = append(re.Words, n.Text)
		}
	}
	if len(re.Words) == 0 {
		if re.Merge {
			re.MergeT()
		} else {
			if re.Xml == "" {
				re.GetXml()
			}
			getW()
		}
	}
	re.Words = Datas.RemoveRepByMapString(re.Words)
	return re.Words
}

// 获取 xml
func (re *RpExcel) GetXml() {
	if err := GetOfficeDocumentCB(re.InFile, re.XmlName, func(s string) {
		re.Xml = strings.Replace(s, "xml:space", "xml_space", -1)
	}); err != nil {
		re.Logger.Println("[error]", re.InFile)
		re.Logger.Println("[error]", err.Error())
	}
}

// 合并 w:t 节点
func (re *RpExcel) MergeT() {
	if re.Xml == "" {
		re.GetXml()
	}
	re.Xml = RepSpacePreserve(re.Xml, "t")
	if len(re.Xml) == 0 {
		re.Logger.Println("[error]", re.InFile)
		re.Logger.Println("[error]", "xml length is zero")
		return
	}
	doc := xmldom.Must(xmldom.ParseXML(re.Xml))
	root := doc.Root
	for _, n := range root.Query("//si") {
		s := ""
		eles := n.FindByName("t")
		for _, t := range eles {
			space := t.GetAttribute("xml_space")
			if space != nil {
				if space.Value == "right" {
					s += t.Text + " "
				} else if space.Value == "left" {
					s += " " + t.Text
				} else if space.Value == "both" {
					s += " " + t.Text + " "
				}
				t.SetAttributeValue("xml_space", "preserve")
			} else {
				s += t.Text
			}
		}
		var first *xmldom.Node
		if len(eles) == 1 {
			re.Words = append(re.Words, s)
			continue
		} else {
			for _, r := range n.FindByName("r") {
				if first == nil {
					first = r
				}
				n.RemoveChild(r)
			}
		}
		if first != nil {
			node := first.GetChild("t")
			if node != nil {
				node.Text = s
				re.Words = append(re.Words, s)
			}
			n.AppendChild(first)
		}
	}
	re.Xml = doc.XML()
	// 对提取到的文件进行去重
	re.Words = Datas.RemoveRepByMapString(re.Words)
}

// 去除替换内容的空格
func (re *RpExcel) Trim() {
	for ind, w := range re.Words {
		re.Words[ind] = strings.TrimSpace(w)
	}
}

// 替换内容
func (re *RpExcel) Replace(words map[string]string) {
	if len(re.Xml) == 0 {
		re.GetXml()
	}
	doc := xmldom.Must(xmldom.ParseXML(re.Xml))
	root := doc.Root
	for _, n := range root.Query("//t") {
		if _, ok := words[n.Text]; ok {
			n.Text = words[n.Text]
		}
	}
	re.Xml = root.XML()
}

// 写出文件
func (re *RpExcel) WriteOut() {
	if err := DoActionInOfficeFileXml(
		re.InFile,
		re.OutFile,
		[]string{
			re.XmlName,
			re.SheetXmlName,
		},
		func(fileName, s string) string {
			if fileName == re.XmlName {
				return re.Xml
			} else if fileName == re.SheetXmlName {
				return "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n" + ToDocxXmlX(re.SheetXml, "_", ":", 8)
			}
			return ""
		}); err != nil {
		re.Logger.Println("[error]", re.InFile)
		re.Logger.Println("[error]", err.Error())
	}
}

// 下面是 20200423 添加内容，关于将 sheet name 获取并修改的部分
// 就当作接口兼容，这里另外写一部分吧
func (re *RpExcel) GetSheetXML() {
	if err := GetOfficeDocumentCB(re.InFile, re.SheetXmlName, func(s string) {
		re.SheetXml = FroXmlParseX(s, ":", "_", 8)
	}); err != nil {
		re.Logger.Println("[error]", re.InFile)
		re.Logger.Println("[error]", err.Error())
	}
}

// 获取 sheet 名称
func (re *RpExcel) GetSheetName() []string {
	if len(re.SheetName) == 0 {
		re.GetSheetXML()
	}
	doc := xmldom.Must(xmldom.ParseXML(re.SheetXml))
	root := doc.Root
	sheets := root.Query("//sheet")
	for _, sheet := range sheets {
		re.SheetName = append(re.SheetName, sheet.GetAttributeValue("name"))
	}
	re.SheetName = Datas.RemoveRepByMapString(re.SheetName)
	return re.SheetName
}

// 替换 sheet 名称
func (re *RpExcel) ReplaceSheetName(repWord map[string]string) {
	if len(re.SheetName) == 0 {
		re.GetSheetXML()
	}
	doc := xmldom.Must(xmldom.ParseXML(re.SheetXml))
	root := doc.Root
	sheets := root.Query("//sheet")
	for _, sheet := range sheets {
		name := sheet.GetAttributeValue("name")
		if _, ok := repWord[name]; ok {
			sheet.SetAttributeValue("name", repWord[name])
		}
	}
	re.SheetXml = root.XML()
}

// 将内容用uuid替换
// 这个是生成模板是，将 sheet name 用 uuid 表示，但是 sheet name 长度不能大于 31
func (re *RpExcel) ReplaceSheetNameToUUID(repWord map[string]string) {
	for k, v := range repWord {
		repWord[k] = Datas.UUID16236(v)
	}
	if len(re.SheetName) == 0 {
		re.GetSheetXML()
	}
	doc := xmldom.Must(xmldom.ParseXML(re.SheetXml))
	root := doc.Root
	sheets := root.Query("//sheet")
	for _, sheet := range sheets {
		name := sheet.GetAttributeValue("name")
		if _, ok := repWord[name]; ok {
			sheet.SetAttributeValue("name", repWord[name])
		}
	}
	re.SheetXml = root.XML()
}

// 将内容用uuid替换
// 这个是将模板中的 sheet name 替换成其他内容，因为 上面这个方法的干扰，需要另外写一个方法
func (re *RpExcel) ReplaceSheetNameFromUUID(_repWord map[string]string) {
	repWord := map[string]string{}
	for k, v := range _repWord {
		repWord[Datas.UUID16236(k)] = v
	}
	if len(re.SheetName) == 0 {
		re.GetSheetXML()
	}
	doc := xmldom.Must(xmldom.ParseXML(re.SheetXml))
	root := doc.Root
	sheets := root.Query("//sheet")
	for _, sheet := range sheets {
		name := sheet.GetAttributeValue("name")
		//name = Datas.UUID36216(name)
		if _, ok := repWord[name]; ok {
			sheet.SetAttributeValue("name", repWord[name])
		}
	}
	re.SheetXml = root.XML()
}
