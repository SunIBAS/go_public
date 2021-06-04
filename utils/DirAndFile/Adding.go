package DirAndFile

import (
	b64 "encoding/base64"
	"io/ioutil"
)

func ReadFileAsBase64String(file string) string {
	content, _ := ioutil.ReadFile(file)
	return b64.StdEncoding.EncodeToString(content)
}

func SaveBase64StringToByteFile(base64, filePath string) {
	content, _ := b64.StdEncoding.DecodeString(base64)
	WriteWithIOUtilByte(filePath, content)
}

func Base64ToByte(base64 string) []byte {
	content, _ := b64.StdEncoding.DecodeString(base64)
	return content
}
