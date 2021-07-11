mkdir ./_release
:: 构建 cov19 项目
cd ./cov19
statik -src=views
cd ..
go build -o ./_release/cov19.exe ./cov19/main.go

:: 构建 downloadBtcData
go build -o ./_release/btc_down.exe ./downloadBtcData/main.go

:: 构建 下载文件
go build -o ./_release/fileDownload.exe ./fileDownload/main.go

::
go build -o ./_release/serverUtils.exe ./serverUtils/main.go

go build -o ./_release/translateOffice.exe ./translateOffice/main.go

go build -o ./_release/生成指定位数助记词.exe ./bit/生成指定位数助记词/main.go

go build -o ./_release/暴力.exe ./bit/随机匹配/main.go

go build -o ./_release/dearConv19.exe ./dearConv19/main.go

go build -o ./_release/种植结构第一步.exe ./种植结构/第一步面化样本点/main.go

go build -o ./_release/scurl.exe ./SimpleCurl/main.go

go build -o ./_release/fileServer.exe ./FileServers/main.go
