#!/bin/zsh


FLAG=$1
HEADER='Content-Type:application/ddfs-go'

testFS() {
    echo "===============TEST START==============="
    echo hash🔑: $1
    echo file📄: $2
    url="localhost:8080/blk?hash=$1"
    echo 1. 添加块🧱
    export ALL_PROXY="";  curl  -H $HEADER --data-binary @$2 $url $FLAG

    echo 
    echo 2. 获取块🧱
    export ALL_PROXY="";  curl -# $url -o $2.get $FLAG

    echo 3. 结果对比，EOF说明正确
    cmp $2 $2.get 

    # rm $2.get
    echo "--------------TEST END------------------"
}

testFS test test.blk
testFS doge doge.jpg
