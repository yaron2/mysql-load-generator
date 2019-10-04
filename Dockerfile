FROM golang:1.9  AS builder

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get github.com/go-sql-driver/mysql
RUN go build -o main .

FROM alpine:3.10

COPY --from=builder /app/main /usr/bin/mysql-load-generator
CMD ["/usr/bin/mysql-load-generator"]
