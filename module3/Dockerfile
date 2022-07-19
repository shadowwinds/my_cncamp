# syntax=docker/dockerfile:1
FROM golang:1.17
WORKDIR /app
ENV HTTP_PORT=8080
COPY *.go .
RUN go build -o http-server
CMD ["http-server"]