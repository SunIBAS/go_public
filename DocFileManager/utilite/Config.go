package utilite

import (
	"encoding/json"
	"public.sunibas.cn/go_public/utils/DirAndFile"
	"strings"
)

const ConfigFileName = "config.json"

type Config struct {
	DirDbPath  string `json:"dir_db_path"`
	LangDbPath string `json:"lang_db_path"`
	TmpPath    string `json:"tmp_path"`
	Port       string `json:"port"`
}

func CheckConfigFileExist() bool {
	if e, err := DirAndFile.PathExists(ConfigFileName); err != nil {
		panic(err)
	} else {
		return e
	}
}

func GetConfig() Config {
	var _config Config
	err := json.Unmarshal([]byte(strings.Join(DirAndFile.ReadAsFileAsLine(ConfigFileName), "\n")), &_config)
	if err != nil {
		return _config
	} else {
		panic(err)
	}
}
