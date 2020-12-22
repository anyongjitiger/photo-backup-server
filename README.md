# TaoStorage
### This project is develop by Go language, It is [TaoAlbum](https://github.com/markusleevip/TaoAlbum-android)'s back-end . Storage cellphone album server.
### 本项目是由Go语言开发，是[TaoAlbum](https://github.com/markusleevip/TaoAlbum-android)的后端实现，实现手机相册存储到私有服务器的功能。

## Preconditions 前提条件
###  Install TaoDb 安装TaoDB
项目地址:[TaoDB](https://github.com/markusleevip/taodb)
在output目标已经提供Windows X64平台的可运行版本。
-----------
	go get github.com/markusleevip/taodb
	cd taodb
	./build.sh
	./taodbd -dbPath=/data/taodb -addr=:7398


## Build & Run 编译&运行

-----------
    go get github.com/markusleevip/taostorage
    cd taostorage/main
    go build
    ./main 
		
## Changelog

### Data:7/17/2019
Add the Browse Album feature

### Date: 2020/12/4
由于开发时候将项目运行在WSL里面，IP地址不固定，所有需要经常修改几个地方：
network_security_config.xml 文件下面的：
 <domain includeSubdomains="false">172.24.224.1</domain>
upload.js 文件下面的：
common_url = 'http://172.24.224.1:8000/'

IP地址来源：网络连接 --> vEthernet (vEthernet (ubun)的ip

    



