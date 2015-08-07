build:
	GOPATH=~/go PATH=\"$(PATH)\":~go/bin go build devicefarm-cli.go

install:
	GOPATH=~/go PATH=\"$(PATH)\":~go/bin go get github.com/PuerkitoBio/goquery
	GOPATH=~/go PATH=\"$(PATH)\":~go/bin go get golang.org/x/net/publicsuffix
	GOPATH=~/go PATH=\"$(PATH)\":~go/bin go get github.com/aws/aws-sdk-go/service/devicefarm
