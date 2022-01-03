export GOPATH="/opt/homebrew/Cellar/go/1.17.3"
export PATH="$GOPATH/libexec/bin:$PATH"
export GO111MODULE="on"
export GOPROXY="http://goproxy.io"
go install github.com/dubbogo/tools/cmd/protoc-gen-go-triple@v1.0.5
protoc --go_out=. --go-triple_out=. DemoService.proto