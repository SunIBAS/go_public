mkdir ./_release
:: 构建 cov19 项目
cd ./cov19
statik -src=views
cd ..
go build -o ./_release/cov19.exe ./cov19/main.go