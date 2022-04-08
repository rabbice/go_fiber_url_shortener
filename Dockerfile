# stage 2
FROM golang:alpine as builder

RUN apk add git

WORKDIR /go/src/github.com/rabbice/url_shortener
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# stage 3
FROM alpine

WORKDIR /root/

COPY --from=builder /go/src/github.com/rabbice/url_shortener/app .

EXPOSE 4000

CMD ["./app"]