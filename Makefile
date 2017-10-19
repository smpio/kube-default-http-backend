GOOS?=linux
GOARCH?=amd64

build:
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -a -ldflags "-s -w" -o default-http-backend .
