FROM golang:1.19  AS builder

WORKDIR /go/src/app
ADD . /go/src/app
RUN go mod tidy -v && go build -o main .

FROM debian:bullseye-slim
RUN mkdir /app
COPY --from=builder /go/src/app/main /app/mysql-loader
COPY config.json /app/config.json
CMD ["/app/mysql-loader"]
