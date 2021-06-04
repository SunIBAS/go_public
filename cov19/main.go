package main

import (
	"public.sunibas.cn/go_public/cov19/Config"
	"public.sunibas.cn/go_public/cov19/Server"
	sqlsint "public.sunibas.cn/go_public/cov19/sql"
)

func main() {
	sqlsint := sqlsint.SqlsInt{
		//DbPath:   ".\\cov19\\db\\db.db",
		DbPath:   "C:\\Users\\IBAS\\go\\src\\public.sunibas.cn\\go_public\\cov19\\db\\db.db",
		Database: nil,
		ShowSql:  true,
	}
	sqlsint.Init()
	sc := Config.ServerConfig{
		Port: ":8081",
	}
	Server.Server(sc, sqlsint)
}
