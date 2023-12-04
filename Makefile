
NAME= $(notdir $(shell pwd))
TAG=$(shell git tag)

init.go:
	sh get_extract_give_name.sh


build: | init.go
	@go build -ldflags '-w -s -X main.Version=${NAME}-${TAG}'
	@notify-send 'Build Complete' 'Your project has been build successfully!' -u normal -t 7500 -i checkbox-checked-symbolic

clean:
	 @go clean
	 @rm -fr init.go geoUR*

