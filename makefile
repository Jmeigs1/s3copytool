run:
	go run .
build:
	go build
install:
	go build && go install

release: release_osx
release_osx:
	GOOS=darwin GOARCH=arm64 go build -o dist/s3copytool_darwin_arm64
	GOOS=darwin GOARCH=amd64 go build -o dist/s3copytool_darwin_amd64

clean:
	rm -f s3copytool dist/*