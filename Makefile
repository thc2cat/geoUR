
NAME= $(notdir $(shell pwd))
TAG=$(shell git tag)

update: 
	@sh init.sh

init.go:
	@sh init.sh


build: | init.go
	@go build -ldflags '-w -s -X main.Version=${NAME}-${TAG}'
	@notify-send 'Build Complete' 'Your project has been build successfully!' -u normal -t 7500 -i checkbox-checked-symbolic

clean:
	 @go clean
	 @rm -fr init.go geoUR*

