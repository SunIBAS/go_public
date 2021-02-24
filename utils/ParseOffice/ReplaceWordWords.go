package ParseOffice

import (
	"github.com/subchen/go-xmldom"
	"log"
	"os"
	"public.sunibas.cn/go_public/utils/Datas"
	"strings"
)

type RpWord struct {
	InFile  string
	OutFile string
	XmlName string
	Words   []string
	Xml     string
	Merge   bool
	Logger  *log.Logger
}

// 初始化文档对象
func (re *RpWord) Init() {
	re.XmlName = "word/document.xml"
	re.Xml = ""
	re.Words = []string{}
	if re.Logger == nil {
		re.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}
}

// 获取文档文字信息
func (re *RpWord) GetWords() []string {
	getW := func() {
		doc := xmldom.Must(xmldom.ParseXML(re.Xml))
		root := doc.Root
		for _, n := range root.Query("w_t") {
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
func (re *RpWord) GetXml() {
	if err := GetOfficeDocumentCB(re.InFile, re.XmlName, func(s string) {
		re.Xml = FroDocXmlParse(s)
	}); err != nil {
		re.Logger.Println("[error]", re.InFile)
		re.Logger.Println("[error]", err.Error())
	}
}

// 处理 xml 的结构（配合 xmldom.Must 解析时生成了额外信息，这里需要额外处理）
func dearTNode(t *xmldom.Node) string {
	space := t.GetAttribute("xml_space")
	s := ""
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
	return s
}

// 合并 w:t 节点
func (re *RpWord) MergeT() {
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
	wps := root.Query("//w_p")
	rwps := []*xmldom.Node{}
	for _, n := range wps {
		cwp := n.FindByName("w_p")
		if len(cwp) > 1 {
			//fmt.Println("drop one")
		} else {
			rwps = append(rwps, n)
		}
	}
	for _, n := range rwps {
		{
			eles := n.FindByName("w_t")
			if len(eles) == 0 {
				continue
			}
		}
		s := ""
		var first *xmldom.Node
		wrs := n.FindByName("w_r")
		for _, wr := range wrs {
			t := wr.GetChild("w_t")
			br := wr.GetChild("w_br")
			if br != nil {
				if first != nil {
					first.GetChild("w_t").Text = s
					re.Words = append(re.Words, s)
					s = ""
					n.AppendChild(first)
					first = nil
				}
				if nil != t {
					//first = xmldom.Must(xmldom.ParseXML(wr.XML())).Root
					s += dearTNode(t)
					//br = first.GetChild("w_br")
					//first.RemoveChild(br)
					// 像 w:smartTag 的 w:r 不是在 w:p 底下的

					wr.RemoveChild(t)
				}
				for wr.Parent != n {
					wr = wr.Parent
				}
				n.RemoveChild(wr)
				n.AppendChild(wr)
			} else if t != nil {
				s += dearTNode(t)
				if first == nil {
					first = wr
				}
				for wr.Parent != n {
					wr = wr.Parent
				}
				n.RemoveChild(wr)
			}
		}

		if first != nil {
			first.GetChild("w_t").Text = s
			re.Words = append(re.Words, s)
			s = ""
			n.AppendChild(first)
			first = nil
		}
		//fmt.Println(s)
	}
	re.Xml = doc.XML()
	re.Words = Datas.RemoveRepByMapString(re.Words)
}

// 去除替换内容的空格
func (re *RpWord) Trim() {
	for ind, w := range re.Words {
		re.Words[ind] = strings.TrimSpace(w)
	}
}

// 替换内容
func (re *RpWord) Replace(words map[string]string) {
	if len(re.Xml) == 0 {
		re.GetXml()
	}
	doc := xmldom.Must(xmldom.ParseXML(re.Xml))
	root := doc.Root
	for _, n := range root.Query("//w_p") {
		eles := n.FindByName("w_t")
		// 跳过 len(eles)
		if len(eles) == 0 {
			continue
		}
		for _, t := range eles {
			s := strings.TrimSpace(t.Text)
			if _, ok := words[s]; ok {
				t.Text = words[s]
			}
		}
	}
	re.Xml = root.XML()
}

// 写出文件
func (re *RpWord) WriteOut() {
	if err := DoActionInOfficeFileXml(
		re.InFile,
		re.OutFile,
		[]string{re.XmlName},
		func(fileName, s string) string {
			return ToDocxXml(re.Xml)
		}); err != nil {
		re.Logger.Println("[error]", re.InFile)
		re.Logger.Println("[error]", err.Error())
	}
}
