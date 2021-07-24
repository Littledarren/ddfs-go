


.PHONY: build test


build:
	bash .github/workflows/filter.sh changed.txt
test:
	echo 'not test'

