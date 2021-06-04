package ParseOffice

import (
	"public.sunibas.cn/go_public/utils/Datas"
	"regexp"
	"strconv"
	"strings"
)

// <w:xxx => <w_xxx
// w:xxx= => w_xxx=
func replaceAllDocTag(xml string, tarFrom, tarTo string) string {
	replace := func(reg, xml string, tarFrom, tarTo string) string {
		re, _ := regexp.Compile(reg)
		all := re.FindAllString(xml, -1)
		all = Datas.RemoveRepByMapString(all)
		for _, a := range all {
			from := strings.Index(a, tarFrom)
			t := a[0:from] + tarTo + a[from+1:]
			xml = strings.Replace(xml, a, t, -1)
		}
		return xml
	}
	xml = replace("<[a-zA-Z0-9]+"+tarFrom+"[a-zA-Z0-9]+", xml, tarFrom, tarTo)
	xml = replace("</[a-zA-Z0-9]+"+tarFrom+"[a-zA-Z0-9]+", xml, tarFrom, tarTo)
	xml = replace("[a-zA-Z0-9]+"+tarFrom+"[a-zA-Z0-9]+=\"", xml, tarFrom, tarTo)
	xml = replace("[a-zA-Z0-9]+"+tarFrom+"[a-zA-Z0-9]+=\"", xml, tarFrom, tarTo)
	return xml
}

func RepSpacePreserve(xml, domStr string) string {
	stypeName := []string{
		`<` + domStr + ` xml_space="left">`,
		`<` + domStr + ` xml_space="right">`,
		`<` + domStr + ` xml_space="both">`,
	}
	for {
		// -1 0 左 1 右 2 两边
		stype := -1
		ind := strings.Index(xml, `<`+domStr+` xml_space="preserve">`)
		if ind < 0 {
			break
		}
		ind += len(`<` + domStr + ` xml_space="preserve">`)
		if xml[ind:ind+1] == " " {
			stype = 0
		}
		rind := strings.Index(xml[ind:], `</`+domStr+`>`)
		if xml[ind:][rind-1:rind] == " " {
			if stype == 0 {
				stype = 2
			} else {
				stype = 1
			}
		}
		if stype > -1 {
			xml = strings.Replace(xml, `<`+domStr+` xml_space="preserve">`, stypeName[stype], 1)
		} else {
			xml = strings.Replace(xml, `<`+domStr+` xml_space="preserve">`, `<`+domStr+`>`, 1)
		}
	}
	return xml
}

func FroDocXmlParse(xml string) string {
	return RepSpacePreserve(replaceAllDocTag(xml, ":", "_"), "w_t")
}
func FroXmlParse(xml string) string {
	return replaceAllDocTag(xml, ":", "_")
}

func ToDocxXml(xml string) string {
	//xml = strings.Replace(xml,`xml_space="left"`,`xml_space="preserve" `,-1)
	//xml = strings.Replace(xml,`xml_space="right"`,`xml_space="preserve" `,-1)
	//xml = strings.Replace(xml,`xml_space="both"`,`xml_space="preserve" `,-1)
	return replaceAllDocTag(xml, "_", ":")
}

// <w:xxx => <w_xxx
// w:xxx= => w_xxx=
func replaceAllDocTagX(xml string, tarFrom, tarTo string, maxLen int) string {
	replace := func(reg, xml string, tarFrom, tarTo string) string {
		re, _ := regexp.Compile(reg)
		all := re.FindAllString(xml, -1)
		all = Datas.RemoveRepByMapString(all)
		for _, a := range all {
			from := strings.Index(a, tarFrom)
			t := a[0:from] + tarTo + a[from+1:]
			xml = strings.Replace(xml, a, t, -1)
		}
		return xml
	}
	xml = replace("<[a-zA-Z0-9]{1,"+strconv.Itoa(maxLen)+"}"+tarFrom+"[a-zA-Z0-9]+", xml, tarFrom, tarTo)
	xml = replace("</[a-zA-Z0-9]{1,"+strconv.Itoa(maxLen)+"}"+tarFrom+"[a-zA-Z0-9]+", xml, tarFrom, tarTo)
	xml = replace("[a-zA-Z0-9]{1,"+strconv.Itoa(maxLen)+"}"+tarFrom+"[a-zA-Z0-9]+=\"", xml, tarFrom, tarTo)
	xml = replace("[a-zA-Z0-9]{1,"+strconv.Itoa(maxLen)+"}"+tarFrom+"[a-zA-Z0-9]+=\"", xml, tarFrom, tarTo)
	return xml
}

// 因为 xlsx 的标签是 mc:xxx xl5ac:xxx 冒号前那部分不是一个字符
func FroXmlParseX(xml, tarFrom, tarTo string, maxLen int) string {
	return replaceAllDocTagX(xml, tarFrom, tarTo, maxLen)
}
func ToDocxXmlX(xml, tarFrom, tarTo string, maxLen int) string {
	return replaceAllDocTagX(xml, tarFrom, tarTo, maxLen)
}
