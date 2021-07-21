#!/bin/zsh


HEADER='Content-Type:application/ddfs-go'

testFS() {
    echo test hash: $1
    echo test file: $2
    url="localhost:8080/blk?hash=$1"
    echo 1. 添加块
    export ALL_PROXY="";  curl -# -H $HEADER -X POST --data-binary @$2 $url  -v 

    echo 2. 获取块
    export ALL_PROXY="";  curl -# $url -o $2.get -v; cmp $2 $2.get 
}

testFS test test.blk
testFS doge doge.jpg
