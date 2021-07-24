#!/usr/bin/env bash

# git根目录
ROOT_DIR=$(git rev-parse --show-toplevel)
cd ${ROOT_DIR}|| exit -1

OUTPUT_FILE=${1:-"changed.txt"}

# 1. 找到变化的文件
if [ -f ${OUTPUT_FILE} ]; then
    echo ${OUTPUT_FILE} Already Exist. Will Remove It.
    rm ${OUTPUT_FILE}
fi
git diff --name-status | awk '{print $2}'  > ${OUTPUT_FILE}
# | grep -v -E '*.yml'

# 2. 提取全量目录
find_target_dir() {
    BASE=$1
    TARGET=$2
    DIR=$(dirname ${BASE})
    while [ ${DIR} != "." ]
    do
        if [ -f "${DIR}/${TARGET}" ]; then
            echo ${DIR}
            return 0
        fi
        DIR=$(dirname ${DIR})
    done
    return 1
}

# 3. 找到main.go所在的位置，编译
cat ${OUTPUT_FILE} | while IFS= read -r line
do
    if dir=$(find_target_dir ${line} "main.go"); then
        (cd ${line}||(echo "failed to cd" && exit -1); go mod tidy && go build )
    fi
done
