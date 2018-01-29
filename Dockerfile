FROM golang:alpine3.7
RUN apk update && apk upgrade && \
    apk add --no-cache bash gcc musl-dev make protobuf protobuf-dev git
WORKDIR /go/src/github.com/pondohva/profile

ADD . .
RUN go get -u google.golang.org/grpc && go get -u google.golang.org/grpc/metadata
RUN go get -u github.com/golang/protobuf/protoc-gen-go && \
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway && \
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
RUN make proto
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o profile

FROM alpine:3.6
RUN apk add --update ca-certificates && \
    rm -rf /var/cache/apk/* /tmp/*
WORKDIR /app
COPY --from=0 /go/src/github.com/pondohva/profile/profile .
RUN apk add --no-cache tini

ENTRYPOINT [ "/sbin/tini", "--", "/app/profile" ]
