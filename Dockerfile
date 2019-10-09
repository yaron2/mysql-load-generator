FROM golang:1.9-alpine

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN apk add git bash
RUN go get github.com/go-sql-driver/mysql
RUN go build -o main .

CMD ["/app/entrypoint.sh"]
