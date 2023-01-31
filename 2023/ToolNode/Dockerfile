FROM golang:1.19-alpine

WORKDIR $GOPATH/src/github.com/alwashali/ToolNode
COPY . .

RUN go mod download

COPY *.go /root/go/src/github.com/alwashali/ToolNode

RUN go build -o /toolnode
EXPOSE 8080

CMD [ "/toolnode" ]