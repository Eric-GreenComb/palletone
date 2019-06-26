# palletone
palletone

## 说明

1. 签名说明，见/script/sign.go

2. /handler/handler_sample.go里是palletone里生成raw tx的函数说明

## 本地运行

go run *.go

## build

1、Mac下编译Linux, Windows平台的64位可执行程序：

$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o palletone .

$ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o palletone .

2、Linux下编译Mac, Windows平台的64位可执行程序：

$ CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o palletone .

$ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o palletone .

3、Windows下编译Mac, Linux平台的64位可执行程序：

$ SET CGO_ENABLED=0SET GOOS=darwin3 SET GOARCH=amd64 go build -o palletone .

$ SET CGO_ENABLED=0 SET GOOS=linux SET GOARCH=amd64 go build -o palletone .

## 发布

生成的palletone是可执行文件，可以写个sh，用pm2管理进程



