all: test

clean:
	rm -f get-started

install: prepare
	godep go install

prepare:
	go get github.com/tools/godep

build: prepare
	godep go build

test: prepare build
	echo "no tests"

deploy:
	ibmcloud cf push

.PHONY: install prepare build test
