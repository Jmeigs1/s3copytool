GO_SOURCES = $(shell find . -type f -name '*.go')

build: s3copytool

run:
	go run .

install:
	go build && go install

s3copytool: $(GO_SOURCES)
	go build

release: release_osx
release_osx:
	GOOS=darwin GOARCH=arm64 go build -o dist/s3copytool_darwin_arm64
	GOOS=darwin GOARCH=amd64 go build -o dist/s3copytool_darwin_amd64

clean:
	rm -f s3copytool dist/*
