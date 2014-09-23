all: build

build:
	go build .

install: build
	cp -f git-fit /usr/local/bin/git-fit

deps:
	go get -u github.com/mitchellh/goamz
