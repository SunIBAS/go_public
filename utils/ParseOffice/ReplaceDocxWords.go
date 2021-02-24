package ParseOffice

import (
	"archive/zip"
	"bytes"
	"os"
	"public.sunibas.cn/go_public/utils/Datas"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"strings"
)

func ReplaceDocxWords(oldFilename, newFilename string, oldWord, newWord []string) (err error) {
	zipReader, err := zip.OpenReader(oldFilename)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	newFile, err := os.Create(newFilename)
	if err != nil {
		return
	}
	defer newFile.Close()

	zipWriter := zip.NewWriter(newFile)
	for _, file := range zipReader.File {
		writer, err := zipWriter.Create(file.Name)
		if err != nil {
			return err
		}
		readCloser, err := file.Open()
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		buf.ReadFrom(readCloser)
		if file.Name == "word/document.xml" {
			allCount := len(oldWord)
			newContent := string(buf.Bytes())
			for i := 0; i < allCount; i++ {
				newContent = strings.Replace(newContent, oldWord[i], newWord[i], -1)
			}
			writer.Write([]byte(newContent))
		} else {
			writer.Write(buf.Bytes())
		}
	}
	zipWriter.Close()
	return nil
}

func DoActionInDocxXml(oldFilename, newFilename string, xmlName []string, Action func(string, string) string) (err error) {
	xmlNames := []string{}
	for _, xml := range xmlName {
		xmlNames = append(xmlNames, "word/"+xml+".xml")
	}
	return DoActionInOfficeFileXml(oldFilename, newFilename, xmlNames, Action)
}

func GetDocument(oldFilename string, outFile string) (err error) {
	zipReader, err := zip.OpenReader(oldFilename)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		readCloser, err := file.Open()
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		buf.ReadFrom(readCloser)

		if file.Name == "word/document.xml" {
			DirAndFile.WriteWithBufio(outFile, string(buf.Bytes()))
		}
	}
	return nil
}

func GetDocumentCB(oldFilename string, cb func(string)) (err error) {
	zipReader, err := zip.OpenReader(oldFilename)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		readCloser, err := file.Open()
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		buf.ReadFrom(readCloser)

		if file.Name == "word/document.xml" {
			cb(string(buf.Bytes()))
		}
	}
	return nil
}

func GetOfficeDocumentCB(oldFilename, xmlName string, cb func(string)) (err error) {
	zipReader, err := zip.OpenReader(oldFilename)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		readCloser, err := file.Open()
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		buf.ReadFrom(readCloser)

		if file.Name == xmlName {
			cb(string(buf.Bytes()))
		}
	}
	return nil
}

func DoActionInOfficeFileXml(oldFilename, newFilename string, xmlNames []string, Action func(string, string) string) (err error) {
	//xmlNames := []string{}
	//for _,xml := range xmlName {
	//	xmlNames = append(xmlNames,"word/" + xml + ".xml")
	//}
	zipReader, err := zip.OpenReader(oldFilename)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	newFile, err := os.Create(newFilename)
	if err != nil {
		return
	}
	defer newFile.Close()

	zipWriter := zip.NewWriter(newFile)
	for _, file := range zipReader.File {
		writer, err := zipWriter.Create(file.Name)
		if err != nil {
			return err
		}
		readCloser, err := file.Open()
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		buf.ReadFrom(readCloser)
		if Datas.ArrayContainItemString(xmlNames, file.Name) {
			newContent := Action(file.Name, string(buf.Bytes()))
			writer.Write([]byte(newContent))
		} else {
			writer.Write(buf.Bytes())
		}
	}
	zipWriter.Close()
	return nil
}
