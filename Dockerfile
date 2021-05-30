FROM golang:1.16 AS builder
COPY . $GOPATH/src/github.com/kentnsw/artical-api
WORKDIR $GOPATH/src/github.com/kentnsw/artical-api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/artapi

FROM scratch
EXPOSE 8080
COPY --from=builder /go/bin/artapi .
ENTRYPOINT ["./artapi"]
