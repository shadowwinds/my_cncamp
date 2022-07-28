FROM golang:1.17 as builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY . .
RUN ls && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM scratch
WORKDIR /app
COPY --from=builder /app/http-server .
ENV HTTP_PORT=8080
EXPOSE ${HTTP_PORT}
ENTRYPOINT  ["./http-server"]