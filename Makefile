build:
	GOPATH=~/go PATH=\"$(PATH)\":~/go/bin go build devicefarm-cli.go

install:
	GOPATH=~/go PATH=\"$(PATH)\":~/go/bin go get github.com/codegangsta/cli
	GOPATH=~/go PATH=\"$(PATH)\":~/go/bin go get github.com/aws/aws-sdk-go/service/devicefarm

gox:
	PATH=\"$(PATH)\":~/go/bin gox -output dist
