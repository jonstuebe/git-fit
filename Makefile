all: build

build:
	go build .

install: build
	cp -f git-fit /usr/local/bin/git-fit

deps:
	go get -u github.com/mitchellh/goamz/aws
	go get -u github.com/mitchellh/goamz/s3
	go get -u github.com/cheggaaa/pb

unittests:
	go test github.com/dailymuse/git-fit/transport github.com/dailymuse/git-fit/util

integrationtests: install
	./integration.py; rm -rf integration

test: unittests integrationtests
