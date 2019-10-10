FROM golang:1.13.1
WORKDIR /go/src/github.com/Calvinn097/pmb-grpc-server
COPY . .
RUN ls
RUN go get -d -v ./...
RUN go install -v ./...

RUN apt-get update
RUN apt-get install -y supervisor
COPY ./docker/supervisord.conf /etc/supervisord.conf

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]