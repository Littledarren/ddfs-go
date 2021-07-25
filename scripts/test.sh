#!/usr/env bash


FLAG=${1}
HEADER='Content-Type:application/ddfs-go'

testBlk() {
    echo "===============TEST START==============="
    echo "hash🔑: ${1}"
    echo "file📄: ${2}"
    url="localhost:8080/blk?hash=${1}"
    echo "1. 添加块🧱"
    export ALL_PROXY="";  curl  -H "${HEADER}" --data-binary @"${2}" "${url}" ${FLAG}

    echo "2. 获取块🧱"
    export ALL_PROXY="";  curl -# "${url}" -o "${2}.get" ${FLAG}

    echo "3. 结果对比，EOF说明正确"
    cmp "${2}" "${2}.get" 

    echo "4. 删除块🧱"
    export ALL_PROXY="";  curl -# -X DELETE "${url}" ${FLAG}

    echo "5. 获取块🧱"
    export ALL_PROXY="";  curl -# "${url}" -v

    # rm ${2}.get
    printf "\\n--------------TEST END------------------\\n"
}


# testBlk test test.blk
# testBlk doge doge.jpg
