FROM golang:1.16-alpine  AS builder

#RUN mkdir /app
#ADD ./main.go /app/
#WORKDIR /app
WORKDIR /go/src/app
ADD ./main.go /go/src/app
RUN apk update && apk add git && go mod init && go get github.com/go-sql-driver/mysql && go get github.com/namsral/flag
RUN go build -o main .

FROM alpine:3.10
RUN mkdir /app
COPY --from=builder /go/src/app/main /app/mysql-loader
COPY config.json /app/config.json
CMD ["/app/mysql-loader"]
