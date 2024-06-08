build:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.version=${VERSION}" -o ${BUILD_OUT} .

build-arm:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s -X main.version=${VERSION}" -o ${BUILD_OUT} .