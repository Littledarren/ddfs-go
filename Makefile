


.PHONY: build test

run: data
	echo "开始执行"
	(cd app/storage/ || echo "请在项目根目录执行" && go mod tidy && go run . & )
	(cd app/tracker/ || echo "请在项目根目录执行" && go mod tidy && go run . &)
build:
	bash .github/workflows/filter.sh changed.txt
test:
	echo 'not test'
data:
	if [ ! -d app/storage/data ]; then mkdir app/storage/data || echo "无法创建data文件夹"; fi;


