FROM golang:1.8
WORKDIR /go/src/github.com/Calvin097/pmb-grpc-server
COPY . .
RUN pwd
RUN ls
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["account_server"]